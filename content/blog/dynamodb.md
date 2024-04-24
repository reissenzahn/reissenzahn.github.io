+++
title = "DynamoDB"
date = "2024-04-12"
tags = ["aws"]
subtitle = "Some notes from the The DynamoDB Book"
+++

<!--
Done: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 16
In-progress: 19
Omitted: 12, 15, 22
Remaining: 18, 19, 20, 21
-->

## Overview
- DynamoDB is a fully-managed NoSQL database.
- Two data models are supported:
  - Key-value store: distributed hash table; allows for storing and retrieving records by primary key.
  - Wide-column store: hash table where the value of each record is a B-tree; allows for performing contiguous range queries.
- It provides fast, consistent performance that does not degrade as the amount of stored data increases.
- Operations are performed by making HTTP requests to the DynamoDB API and this HTTP connection model supports a virtually unlimited number of concurrent requests.
- A flexible pricing model is available:
  - Provisioned capacity: read and write throughput can be changed independently and dynamically.
  - On-demand: pay per request.
- DynamoDB is suited for most online, transactional processing (OLTP) applications with high volumes of small bits of data being read and written.
- The rigidity of DynamoDB operations prevents writing inefficient queries and producing data models that do not scale.

## Concepts

### Tables, Items & Attributes
- A table is a collection of items, each composed of one or more attributes.
- An attribute is given a type when it is written:
  - *Scalar*: single simple value (string, number, binary, boolean, null).
  - *Complex*: groupings with arbitrary nested attributes (list, map).
  - *Set*: multiple unique values of the same type (string sets, number sets, binary sets).
- The attribute type affects which operations can be performed.
- Attributes with the same underlying value by a different type are not considered equal.
- A table is schemaless and attributes are not required on every item.
- Multiple different types of entities are commonly stored in a single table to allow for handling complex access patterns in a single request without the need for joins.
- Use prefixes to differentiate between different entity types (e.g. `CUSTOMER#123`).

![items.png](/img/dynamodb/items.png)

### Primary Keys
- A primary key must be declared when creating a table.
- Each item must include the primary key and is uniquely identified by that primary key.
- There are two types of primary keys:
  - *Simple*: consists of a single partition key.
  - *Composite*: consists of a partition key and sort key.
- The terms *hash key* and *range key* are also used to refer to the partition key and sort key, respectively.
- A simple primary key allows for only fetching a single item at a time while a composite primary key enables fetching all items with the same partition key.
- Writing an item using a primary key that already exits will overwrite the existing item (unless this is explictly disabled, in which case the write will be rejected).

### Secondary Indexes
- Secondary indexes allow for reshaping data into another format for querying.
- All items will be copied from the base table into the secondary index in the reshaped form.
- A key schema consisting of a partition key and (optional) sort key must be declared when creating a secondary index.
- There are two types of secondary indexes:
  - *Local*: uses the same partition key as the primary key but a different sort key.
  - *Global*: uses any attributes for the partition key and sort key.
- Filtering is built into the data model as the primary keys and secondary indexes determine how data is retrieved.
- Separate application attributes (*FirstName*, *LastUdated*, etc.) from indexing attributes (*PK*, *GSI1SK*, etc.) and avoid reusing attribute across multiple indexes to avoid unnecessary complexity.

| | LSI | GSI |
|-|-|-|
| Throughput | Shared with base table | Separately provisioned |
| Consistency | Allow opting into strongly-consistent reads | Eventually-consistent reads only |
| Creation time | Must be specified when table created | Can be created and deleted as needed |

### Item Collections
- An item collection consists of the group of items that share the same partition key in either the base table or a secondary index.
- All the items in an item collection will be allocated to the same partition.
- An item collection is ordered and stored as a B-tree.
- All items in an LSI are included as part of the same item collection as the base table.
- The item collection size includes both the size of the items in the base table and the size of the items in the local secondary index.

![item-collections.png](/img/dynamodb/item-collections.png)

### Attribute Projections
- A projection is the set of attributes that is copied from a table into a secondary index.
- The partition key and sort key of the table are always projected into the index and other attributes can be projected as required.
- This must be declared when the secondary index is created. There are three options:
  - *KEYS_ONLY*: project the table partition key and sort key values, plus the index key values.
  - *INCLUDE*: Include other additional non-key attributes.
  - *ALL*: The secondary index includes all of the attributes from the source table.
- Choosing which attributes to project presents a trade-off between throughput costs and storage costs.
- It is not possible to change projected attributes once the index has been created.

### Sparse Indexes
- When an item is written to the base table, it will be copied into the secondary index only if it has the elements of the key schema for the secondary index.
- If the item does not have those elements then it will not be copied into the secondary index.
- A sparse index is a secondary index that intentionally excludes certain items from the base table to help satisfy certain access patterns.
- An attribute can be removed from an item to remove it from a sparse index.

### Partitioning
- Data is sharded across multiple partitions based on the partition key which allows for horizontal scaling by adding more storage nodes.
- There are three storage nodes for each partition: a primary node that holds the canonical data and two secondary nodes that provide durability and serve reads.
- The request router serves as the frontend for all requests.
- An incoming write request is handled as follows:
  1. The request router authenticates the request.
  2. The request router uses the hash of the partition key to route the request to the appropriate primary node.
  3. The primary node commits the write and also commits the write to one of the two secondary nodes.
  4. The primary node responds to the client to indicate that the write was successful.
  5. The primary node asynchronously replicates the write to the third storage node.
- Adaptive capacity automatically spreads throughput around a table to the items that require it.

### Consistency
- Consistency refers to whether a particular read operation receives all prior write operations.
- There are two consistency options available:
  - *Strong*: Any item read will reflect all prior writes.
  - *Eventual*: It is possible that read items will not reflect all prior writes.
- Reads are eventually-consistent by default thought it is possible to opt into strongly-consistent for base tables and LSIs.
- An eventually-consistent read consumes half the read capacity of a strongly-consistent read.
- Sources of eventual consistency include asynchronous data replication from primary to secondary nodes and from base tables to GSIs.

### Key Overloading 
- Key overloading refers to using generic names for the primary keys and using different values dependending on the item type.
- Prefixes are used for the partition and sort key values in order to identify the item type and to avoid overlap between different item types.
- Secondary indexes can be overloaded just like primary keys.
- Common generic names are: *PK*, *SK*, *GSI1PK*, *GSI1SK*, etc.

![key-overloading.png](/img/dynamodb/key-overloading.png)

### Time-to-Live
- TTLs allow for automatically deleting items after a specified time on a per-item basis.
- To use TTLs, an attribute is specified on the table that will serve as the marker for item deletion.
- The attribute must be of type number.
- To expire an item, a Unix timestamp at seconds granularity can be stored in the specified attribute that indicates that time after which the item should be deleted.
- Items are usually deleted within 48 hours after the time indicated by the attribute.

### Capacity Units
- A RCU represents one strongly-consistent read per second or two eventually-consistent reads per second for an item up to 4KB. Transactional read requests require two RCUs to perform one read per second for items up to 4KB. Reading an item larger than 4KB will consume additional RCUs.
- One WCU represnts one write per second for an item up to 1KB in size. Transactional write requests require two WCUs to perform one write per second for items up to 1KB. Writing an item larger than 1KB will consume additional WCUs.

### Limits
- A single item is limited to 400KB.
- The *Query* and *Scan* actions will read a maximum of 1MB befor paginating (applied before filter and projection expressions).
- A single partition can have a maximum of 3000 RCUs or 1000 WCUs.
- When using a LSI, a single item collection cannot be larger than 10GB.
- There is no limit to the number of items that can be stored in a table.
- Up to 1MB can be returned in a single response.

### Efficiency
- By specifying the partition key, an operation starts with an *O(1)* lookup that reduces the dataset down to a maximum of 10GB on a particular storage node.
- For a *Query* on an item collection of size *n*, an *O(log n)* operation is required to find the starting value for the sort key after which a sequential read is limited to 1MB.

### Sorting
- Items need to be arranged so that they are sorted in advance.
- If a specific ordering is required when retrieving multiple items then a composite primary key must be used such that the ordering is performed with the sort key.
- Sort keys of type string or binary are sorted in order of UTF-8 bytes.
- Given uppercase letters are sorted before lowercase letters, sort keys should be standardized to a single case to avoid unexpected behavior.
- For timestamps to be sortable they should use a sortable format like a Unix epoch or ISO-8601.

### Streams
- Streams provide support for change data capture.
- Whenever an item is written, updated or deleted, a record containing the details of that mutation will be written to the stream.

## Operations

### Item-based
- Item-based actions operate on specific items: *GetItem*, *PutItem*, *UpdateItem* and *DeleteItem*.
- These are the only actions that can be used to write, update or delete items.
- The full primary key must be specified and the operation must be performed on the base table.
- *PutItem* can overwrite an existing item with the same primary key.
- *UpdateItem* will create an item if it does not exist, otherwise it will only alter the properties specified.

### Batch/Transaction
- Batch and transaction actions are used for operating on multiple items in a single request.
- These actions must specify the full primary key of items.
- Batch actions allow for reads or writes to succeed or fail independently.
- With transaction actions, the failure of a single operation will cause all the other writes to be rolled back.
- The *TransactWriteItem* action allows for including up to 10 items in a single request which can be a combination of the following: *PutItem*, *UpdateItem*, *DeleteItem* and *ConditionCheck*.

### Query
- The *Query* action can be used to retrieve a contiguous block of items within a single item collection.
- The partition key must be specified.
- Various conditions can be specified on the sort key: *>=*, *<=*, *begins_with()*, or *BETWEEN*.
- This can be performed against the base table or a secondary index.

### Scan
- A *Scan* will retrieve all items in a table.
- In exceptional situations, a sparse secondary index can be modeled in a way that expects a *Scan*.
- Two optional properties are available to enable parallel scans:
  - *TotalSegments*: The total number of segments to split the *Scan* across.
  - *Segment*: The segment number to be scanned by this particular request.

### Optional Properties
- *ConsistentRead* is used to opt into strongly-consistent reads using *GetItem*, *BatchGetItem*, *Query* and *Scan*.
- *ScanIndexForward* controls which way a *Query* will read results from the sort key. By default, ascending order is used; using *ScanIndexForward=False* will use descending order.
- *ReturnValues* determines which values a *PutItem*, *UpdateItem*, *DeleteItem* or *TransactWriteItem* operation returns:
  - *ALL_OLD*: Return all the attributes from the item before the operation was applied.
  - *UPDATED_OLD*: For any attributes updated in the operation, return the attributes before the operation was applied.
  - *ALL_NEW*: Return all the attributes from the item after the operation was applied.
  - *UPDATED_NEW*: For any attributes updated in the operation, return the attributes after the operation was applied.
- *ReturnConsumedCapacity* returns data about the capacity units that were consumed by the request.
- *ReturnItemCollectionMetrics* returns item collection size metrics.

![scan-index-forward.png](/img/dynamodb/scan-index-forward.png)

## Expressions

### Placeholders
- There are two types of placeholders:
  - Expression attribute values (e.g. `:foo`) represent an attribute value being evaluated in the request.
  - Expression attribute names (e.g. `#bar`) specify the names of attributes.
- Expression attribute values allow for declaring attribute types without complicating expression parsing.
- It is not required to use expression attribute names though they are required if an attribute name conflicts with a reserved word.

### Key Condition
- A key condition expression is used in a *Query* to describe which items to retrieve.
- It can only reference elements of the primary key and must include the partition key.
- Sort key conditions can use simple comparisons: *>*, *<*, *=*, *begins_with*, *BETWEEN*.
- Every condition on the sort key can be expressed with the *BETWEEN* operator.
```py
results = client.query(
  TableName='CustomerOrders',
  KeyConditionExpression='#c = :c AND #ot BETWEEN :start AND :end',
  ExpressionAttributeNames={
    '#c': 'CustomerId',
    '#ot': 'OrderTime'
  },
  ExpressionAttributeValues={
    ':c': { 'S': '36ab55a589e4' },
    ':start': { 'S': '2020-01-10T00:00:00.000000' },
    ':end': { 'S': '2020-01-20T00:00:00.000000' }
  }
)
```

### Filter
- A filter expression can be used in *Query* and *Scan* operations to describe which items should be returned to the client after finding the items that match the key conditions.
- A filter expression can be applied on any attribute in the table.
- *Query* and *Scan* operations will return a maximum of 1MB and this limit is applied before the filter expression.
- This is useful for reducing payload sizes and to filter out items with expired TTLs.
```py
results = client.query(
  TableName='MovieRoles',
  KeyConditionExpression='#actor = :actor',
  FilterExpression='#genre = :genre'
  ExpressionAttributeNames={
    '#actor': 'Actor',
    '#genre': 'Genre'
  },
  ExpressionAttributeValues={
    ':actor': { 'S': 'Joe Bloggs' },
    ':genre': { 'S': 'Drama' }
  }
)
```

### Projection
- A project expression can be used in read operations to describe which attributes to return on read items.
- Projection expressions are evaluated after the items are read from the table and the 1MB limit is reached.
- This can be used to access nested properties in a list or map attribute.
- This can be useful to reduce payload sizes.
```py
result = client.query(
  TableName='MovieRoles',
  KeyConditionExpression='#actor = :actor',
  ProjectionExpression='#actor, #movie, #role, #year, #genre',
  ExpressionAttributeNames={
    '#actor': 'Actor',
    '#movie': 'Movie',
    '#role': 'Role',
    '#year': 'Year',
    '#genre': 'Genre'
  },
  ExpressionAttributeValues={
    ':actor': { 'S': 'Tom Hanks' }
  }
)
```

### Condition
- A condition expression is used in write operation to assert the existing condition (or non-condition) of an item before writing to it.
- The operation will be rejected if the assertion fails.
- This avoids the need for additional requests to fetch an item before manipulating it (and the need for handling race conditions).
- Condition expressions can operate on any attribute on the item.
- Several operators can be used: *>*, *<*, *=*, *BETWEEN*, *attribute_exists()*, *attribute_not_exists()*, *attribute_type()*, *begins_with()*, *contains()* and *size()*.
- The *TransactWriteItem* action can specify combination of different write operations (PutItem, UpdateItem or DeleteItem) or ConditionChecks.
```py
result = dynamodb.put_item(
  TableName='Users',
  Item={
    'Username': { 'S': 'joebloggs123' },
    'Name': { 'S': 'Joe Bloggs' },
    'CreatedAt': { 'S': datetime.datetime.now().isoformat() },
  },
  ConditionExpression='attribute_not_exists(#username)',
  ExpressionAttributeNames={
    '#username': 'Username',
  }
)

result = dynamodb.transact_write_items(
  TransactItems=[
    {
      'ConditionCheck': {
        'Key': {
          'PK': { 'S': 'Admins#<orgId>' }
        },
        'TableName': 'SaasApp',
        'ConditionExpression': 'contains(#a, :user)',
        'ExpressionAttributeNames': {
          '#a': 'Admins'
        },
        'ExpressionAttributeValues': {
          ':user': { 'S': '<username>' }
        }
      }
    },
    {
      'Delete': {
        'Key': {
          'PK': { 'S': 'Billing#<orgId>' }
        },
        'TableName': 'SaasApp'
      }
    }
  ]
)
```

### Update
- An update expression is used in *UpdateItem* to describe the desired updates.
- There are four verbs for stating the changes:
  - *SET*: Add or overwrite an attribute, add or subtract from a number attribute.
  - *REMOVE*: Delete an attribute, delete nested properties from a list or map.
  - *ADD*: Add to a number attribute, insert an element into a set attribute.
  - *DELETE*: Remove an element from a set attribute.
- Any combination of these verbs may be used in a single update statement and multiple operations can be performed for a single verb (e.g. `SET Name = :name, UpdatedAt = :updatedAt REMOVE InProgress`).
- Update expressions can act directly on nested properties within lists and maps.
```py
dynamodb.update_item(
  TableName='Users',
  Key={
    'Username': { 'S': 'joebloggs123' }
  },
  UpdateExpression='SET #phone.#mobile :cell',
  ExpressionAttributeNames={
    '#phone': 'PhoneNumbers',
    '#mobile': 'MobileNumber'
  },
  ExpressionAttributeValues={
    ':cell': { 'S': '+1-555-555-5555' }
  }
)
```

## Patterns

### One-to-Many

#### Denormalization using a complex attribute
- Use an attribute with a complex data type like a list or map.
- It will not be possible to support access patterns based on the values in the complex attribute because a complex attribute cannot be used in a primary key.
- The amount of data in the complex attribute cannot be unbounded.

![one-to-many-1.png](/img/dynamodb/one-to-many-1.png)

#### Denormalization by duplicating data
- Duplicate the parent fields in each of the child items.
- It may be acceptable to duplicate the data depending on how often it changes and how many items contain the duplicated data.
- This balances the benefit of duplication (in the form of faster reads) against the costs of updating the data.

![one-to-many-2.png](/img/dynamodb/one-to-many-2.png)

#### Query with composite primary key
- Use a composite primary key and a *Query* to fetch multiple items within a single item collection.
- This solves for four common access patterns:
  1. Retrieve the parent item using *GetItem* with the primary key of the parent item.
  2. Retrieve the parent item and all its child items using *Query* and the partition key.
  3. Retrieve only the child items using *Query* with `begins_with(SK, "CHILD#")`.
  4. Retrieve a specific child item using *GetItem* with the primary key of the child item.

![one-to-many-3.png](/img/dynamodb/one-to-many-3.png)

#### Query with secondary index
- A *Query* can also be used with a secondary index if the primary key is reserved for some other purpose.
- This could be due to storing hierarchical data with a number of levels.

![one-to-many-4.png](/img/dynamodb/one-to-many-4.png)

![one-to-many-5.png](/img/dynamodb/one-to-many-5.png)

#### Composite sort keys with hierarchical data
- It is not feasible to keep adding secondary indexes to enable arbitrary levels of fetching throughout a data hierarchy.
- A composite sort key refers to combining multiple properties together in the sort key to enable different search granularity.
- With each level separated by a hash in the sort key, we can search at different levels of granularity using `begins_with()`: `PK = <Country> AND begins_with(SK, '<State>#')`, `PK = <Country> AND begins_with(SK, '<State>#<City>')`, `PK = <Country> AND begins_with(SK, '<State>#<City>#<ZipCode>')`, etc.
- This works well when there are more than two levels of hierarchy and access patterns for different levels within the hierarchy and when we want to return all sub-items in a level of the hierarchy.

![one-to-many-6.png](/img/dynamodb/one-to-many-6.png)

### Filtering

#### Assembling different collections of items
- The sort key can enable filtering based on how the items are arranged within a particular item collection even if it does not have inherent meaning.

![filtering-1.png](/img/dynamodb/filtering-1.png)
```py
result = client.query(
  TableName='GitHubTable',
  KeyConditionExpression='#pk = :pk AND #sk <= :sk',
  ExpressionAttributeNames={
    '#pk': 'PK',
    '#sk': 'SK'
  },
  ExpressionAttributeValues={
    ':pk': { 'S': 'REPO#alexdebrie#dynamodb-book' },
    ':sk': { 'S': 'REPO#alexdebrie#dynamodb-book' }
  },
  ScanIndexForward=True
)
```

#### Composite sort key
- A composite sort key contains two or more data elements.
- A composite key can be used as the sort key for a secondary index to allow for querying for an exact match on the first part of the composite key followed by more fine-grained filtering on the remaining value.
- This works well when an access pattern always requires filtering on two or more attributes and one of the attributes is an enum-like value.

![filtering-2.png](/img/dynamodb/filtering-2.png)

![filtering-3.png](/img/dynamodb/filtering-3.png)

```py
result = client.query(
  TableName='CustomerOrders',
  IndexName='OrderStatusDateGSI',
  KeyConditionExpression='#c = :c AND #osd BETWEEN :start and :end',
  ExpressionAttributeNames={
    '#c': 'CustomerId',
    '#ot': 'OrderStatusDate'
  },
  ExpressionAttributeValues={
    ':c': { 'S': '2b5a41c0' },
    ':start': { 'S': 'CANCELLED#2019-07-01T00:00:00.000000' },
    ':end': { 'S': 'CANCELLED#2019-10-01T00:00:00.000000' },
  }
)
```

#### Using sparse indexes to provide a global filter on an item type
- A sparse index can be used to filter within an entity type based on a particular condition by adding an attribute only to those items which fulfil the condition.

![filtering-4.png](/img/dynamodb/filtering-4.png)

![filtering-5.png](/img/dynamodb/filtering-5.png)

#### Using sparse indexes to project a single type of entity
- A single entity type can be projected into a secondary index by only populating a particular attribute on items of that entity type.
- A *Scan* against the secondary index can be used to obtain all the items of that particular entity type.
- This does not work with index overloading as it relies on projecting only a single entity type into the secondary index.

![filtering-6.png](/img/dynamodb/filtering-6.png)

![filtering-7.png](/img/dynamodb/filtering-7.png)

### Sorting

#### Sorting on changing attributes
- Using a frequently changing value as the sort key is not feasible as it would require repeatedly deleting and re-creating the item as it is not possible to update elements of the primary key.
- Instead, immutable values should be used for the primary key and the frequently changing value can be used as the sort key of a secondary index.

![sorting-1.png](/img/dynamodb/sorting-1.png)

#### Sorting with multiple entity types
- When co-locating items for one-to-many relationships, it is necessary to consider the order in which related items need to be returned to determine where to place the parent item.
- This reasoning also applies if there are two relational access patterns in a single item collection in which case one access pattern needs to fetch items in ascending order while the other fetches items in descending order.

![sorting-2.png](/img/dynamodb/sorting-2.png)

![sorting-3.png](/img/dynamodb/sorting-3.png)

```py
result = dynamodb.query(
  TableName='SaaSTable',
  KeyConditionExpression='#pk = :pk AND #sk <= :sk',
  ExpressionAttributeNames={
    '#pk': 'PK',
    '#sk': 'SK'
  },
  ExpressionAttributeValues={
    ':pk': { 'S': 'ORG#MCDONALDS' },
    ':sk': { 'S': 'ORG#MCDONALDS' }
  },
  ScanIndexForward=False
)
```

#### Zero-padding with numbers
- It is sometimes necessary to order items numerically even when the sort key type is a string, such as when using item type prefixes.
- Lexicographic sorting with numbers in strings evaluates one character at a time and so "10" preceeds "2".
- This can be avoided by zero-padding the numbers like "00010" and "00002".

### Other

#### Ensuring uniqueness on two or more attributes
- To ensure that a particular attribute is unique it must be built directly into the primary key structure.
- To ensure that two attributes are unique across the table, it will be necessary to write two items in a transaction: one that tracks the item by the first attribute and another marker item identified by the second attribute.
- Each write operation should include a condition expression that ensures that the written primary key is unique.
```py
response = client.transact_write_items(
  TransactItems=[
    {
      'Put': {
        'TableName': 'UsersTable',
        'Item': {
          'PK': { 'S': 'USER#alexdebrie' },
          'SK': { 'S': 'USER#alexdebrie' },
          'Username': { 'S': 'alexdebrie' },
          'FirstName': { 'S': 'Alex' },
         },
        'ConditionExpression': 'attribute_not_exists(PK)'
      }
    },
    {
      'Put': {
        'TableName': 'UsersTable',
        'Item': {
          'PK': { 'S': 'USEREMAIL#alex@debrie.com' },
          'SK': { 'S': 'USEREMAIL#alex@debrie.com' },
        },
        'ConditionExpression': 'attribute_not_exists(PK)'
      }
    }
  ]
)
```

#### Handling sequential identifiers
- Sequential identifiers can be obtained by using an *UpdateItem* operation to increment a count attribute and return the incremented value.
- This value can then be used to create a new item with that identifier.
```py
result = client.update_item(
  TableName='JiraTable',
  Key={
    'PK': { 'S': 'PROJECT#my-project' },
    'SK': { 'S': 'PROJECT#my-project' },
  },
  UpdateExpression='SET #count = #count + :inc',
  ExpressionAttributeNames={
    '#count': 'IssueCount',
  },
  ExpressionAttributeValues={
    ':incr': { 'N': '1' },
  },
  ReturnValues='UPDATED_NEW'
)

current_count = result['Attributes']['IssueCount']['N']

result = client.put_item(
  TableName='JiraTable',
  Item={
    'PK': { 'S': 'PROJECT#my-project' },
    'SK': { 'S': f"ISSUE#{current_count}" },
    'IssueTitle': { 'S': 'Build DynamoDB data model' }
  }
)
```

#### Paginated responses
- Pagination within an item collection can be achieved by having follow up requests specify the last seen sort key value.

![other-1.png](/img/dynamodb/other-1.png)

```py
# first request
result = client.query(
  TableName='Ecommerce',
  KeyConditionExpression='#pk = :pk, #sk < :sk',
  ExpressionAttributeNames={
    '#pk': 'PK',
    '#sk': 'SK'
  },
  ExpressionAttributeValues={
    ':pk': 'USER#alexdebrie',
    ':sk': 'ORDER$'
  },
  ScanIndexForward=False,
  Limit=5
)

# follow-up request
result = client.query(
  TableName='Ecommerce',
  KeyConditionExpression='#pk = :pk, #sk < :sk',
  ExpressionAttributeNames={
    '#pk': 'PK',
    '#sk': 'SK'
  },
  ExpressionAttributeValues={
    ':pk': 'USER#alexdebrie',
    ':sk': 'ORDER#1YRfXS14inXwIJEf9tO5hWnL2pi'
  },
  ScanIndexForward=False,
  Limit=5
)
```

#### Singleton items
- A singleton item applies across the entire application.

![other-2.png](/img/dynamodb/other-2.png)

#### Reference counts
- A reference count of the number of related child items for a parent item can be maintaned by creating child items in a transaction that also increments a count atttribute on the parent item.
```py
result = dynamodb.transact_write_items(
  TransactItems=[
    {
      'Put': {
        'TableName': 'GitHubModel',
        'Item': {
          'PK': { 'S': 'REPO#alexdebrie#dynamodb-book' },
          'SK': { 'S': 'STAR#danny-developer' },
          # ...
        },
        'ConditionExpression': 'attribute_not_exists(PK)',
      },
    },
    {
      'Update': {
        'TableName': 'GitHubModel',
        'Key': {
          'PK': { 'S': 'REPO#alexdebrie#dynamodb-book' },
          'SK': { 'S': '#REPO#alexdebrie#dynamodb-book' },
        },
        'ConditionExpression': 'attribute_exists(PK)',
        'UpdateExpression': 'SET #count = #count + :inc',
        'ExpressionAttributeNames': {
          '#count': 'StarCount',
        },
        'ExpressionAttributeValues': {
          ':inc': { 'N': '1' },
        },
      },
    },
  ]
)
```

## Examples

1. Identify entities and create an entity-relationship diagram
2. Define the access patterns for each entity and enumerate the index and parameteres used for each one
3. Model the primary key structure for each entity
4. Handle additional access patterns with secondary indexes and streams

### Session Store

#### Entity Relationships

![session-store-1.png](/img/dynamodb/session-store-1.png)

#### Access Patterns
| Access Pattern | Index | Parameters | Notes |
|-|-|-|-|
| Create Session | Main table | SessionToken | Use condition expression to enforce uniqueness |
| Get Session | Main table | SessionToken  | Use filter expression to handle expired tokens |
| Delete Session (time-based) | N/A | N/A | Handled by TTL |
| Delete Sessions for User | UserIndex | Username | Find tokens for user |
|  | Main table | SessionToken | Delete each token |

#### Table Structure
- Simple primary key with a partition key of *SessionToken*.
- GSI called *UserIndex* whose key schema is a simple primary key with a partition key of *Username* and *KEYS_ONLY* projection.
- TTL property set on the attribute named *TTL*.

![session-store-2.png](/img/dynamodb/session-store-2.png)

![session-store-3.png](/img/dynamodb/session-store-3.png)

#### Code Snippets
```py
# prevent duplicate session tokens from being written using a condition expression
created_at = datetime.datetime.now()
expires_at = created_at + datetime.timedelta(days=7)

client.put_item(
  TableName='SessionStore',
  Item={
    'SessionToken': { 'S': str(uuid.uuidv4()) },
    'Username': { 'S': 'dave' },
    'CreatedAt': { 'S': created_at.isoformat() },
    'ExpiresAt': { 'S': expires_at.isoformat() },
  },
  ConditionExpression='attribute_not_exists(SessionToken)',
)

# handle expired items using a filter expression
epoch_seconds = int(time.time())

client.query(
  TableName='SessionStore',
  KeyConditionExpression='#token = :token',
  FilterExpression='#ttl <= :epoch',
  ExpressionAttributeNames={
    '#token': 'SessionToken',
    '#ttl': 'TTL'
  },
  ExpressionAttributeValues={
    ':token': { 'S': '0bc6bdf8-6dac-4212-b11a-81f784297c78' },
    ':epoch': { 'N': str(epoch_seconds) }
  }
)

# query GSI for session tokens related to user and delete each one
results = client.query(
  TableName='SessionStore',
  Index='UserIndex',
  KeyConditionExpression='#username = :username',
  ExpressionAttributeNames={
    '#username': 'Username',
  },
  ExpressionAttributeValues={
    ':username': { 'S': 'alexdebrie' },
  },
)

for result in results['Items']:
  client.delete_item(
    TableName='SessionStore',
    Key={
      'SessionToken': result['SessionToken']
    }
  )
```

### e-Commerce Application

#### Entity Relationships

![e-commerce-1.png](/img/dynamodb/e-commerce-1.png)

| Entity | PK | SK |
|-|-|-|
| Customers | CUSTOMER#<Username> | CUSTOMER#<Username> |




• Create Customer (unique on both username and email address)
• Create / Update / Delete Mailing Address for Customer
• Place Order
• Update Order
• View Customer & Most Recent Orders for Customer
• View Order & Order Items


With our first example out of the way, let’s work on something a
little more complex. In this example, we’re going to model the
ordering system for an e-commerce application. We’ll still have the
training wheels on, but we’ll start to look at important DynamoDB
concepts like primary key overloading and handling relationships
between entities.
Let’s get started.
19.1. Introduction
In this example, you’re working on part of an e-commerce store.
For this part of our application, we’re mostly focused on two core
areas: customer management and order management. A different
service in our application will handle things like inventory, pricing,
and cart management.
I’m going to use screens from ThriftBooks.com to act as our
example application. I love using real-life examples as it means I
don’t have to do any UI design. Fortunately, I have some children
who are voracious readers and a wife that is trying to keep up with
their needs. This means we have some great example data. Lucky
you!
We want to handle the following screens. First, the Account
Overview page:
301
Notice that this page includes information about the user, such as
the user’s email address, as well as a paginated list of the most
recent orders from the user. The information about a particular
order is limited on this screen. It includes the order date, order
number, order status, number of items, and total cost, but it doesn’t
include information about the individual items themselves. Also
notice that there’s something mutable about orders—the order
status—meaning that we’ll need to change things about the order
over time.
To get more information about an order, you need to click on an
order to get to the Order Detail page, shown below.
302
The Order Detail page shows all information about the order,
including the summary information we already saw but also
additional information about each of the items in the order and the
payment and shipping information for the order.
Finally, the Account section also allows customers to save addresses.
The Addresses page looks as follows:
A customer can save multiple addresses for use later on. Each
address has a name ("Home", "Parents' House") to identify it.
303
Let’s add a final requirement straight from our Product Manager.
While a customer is identified by their username, a customer is also
required to give an email address when signing up for an account.
We want both of these to be unique. There cannot be two customers
with the same username, and you cannot sign up for two accounts
with the same email address.
With these needs in mind, let’s build our ERD and list out our
access patterns.
19.2. ERD and Access Patterns
As always, the first step to data modeling is to create our entityrelationship diagram. The ERD for this use case is below:
Our application has four entities with three relationships. First,
304
there are Customers, as identified in the top lefthand corner. A
Customer may have multiple Addresses, so there is a one-to-many
relationship between Customers and Addresses as shown on the
lefthand side.
Moving across the top, a Customer may place multiple Orders over
time (indeed, our financial success depends on it!), so there is a oneto-many relationship between Customers and Orders. Finally, an
Order can contain multiple OrderItems, as a customer may
purchase multiple books in a single order. Thus, there is a one-tomany relationship between Orders and OrderItems.
Now that we have our ERD, let’s create our entity chart and list our
access patterns.
The entity chart looks like this:
Entity PK SK
Customers
Addresses
Orders
OrderItems
Table 16. E-commerce entity chart
And our acccess patterns are the following:
• Create Customer (unique on both username and email address)
• Create / Update / Delete Mailing Address for Customer
• Place Order
• Update Order
• View Customer & Most Recent Orders for Customer
• View Order & Order Items
With these access patterns in mind, let’s start our data modeling.
305
19.3. Data modeling walkthrough
As I start working on a data modeling, I always think about the
same three questions:
1. Should I use a simple or composite primary key?
2. What interesting requirements do I have?
3. Which entity should I start modeling first?
In this example, we’re beyond a simple model that just has one or
two entities. This points toward using a composite primary key.
Further, we have a few 'fetch many' access patterns, which strongly
points toward a composite primary key. We’ll go with that to start.
In terms of interesting requirements, there are two that I noticed:
1. The Customer item needs to be unique on two dimensions:
username and email address.
2. We have a few patterns of "Fetch parent and all related items"
(e.g. Fetch Customer and Orders for Customer). This indicates
we’ll need to "pre-join" our data by locating the parent item in
the same item collection as the related items.
In choosing which entity to start with, I always like to start with a
'core' entity in the application and then work outward as we model
it out. In this application, we have two entities that are pretty
central: Customers and Orders.
I’m going to start with Customers for two reasons:
1. Customers have a few uniqueness requirements, which generally
require modeling in the primary key.
2. Customers are the parent entity for Orders. I usually prefer to
start with parent entities in the primary key.
306
With that in mind, let’s model out our Customer items.
19.3.1. Modeling the Customer entity
Starting with the Customer item, let’s think about our needs around
the Customer.
First, we know that there are two one-to-many relationships with
Customers: Addresses and Orders. Given this, it’s likely we’ll be
making an item collection in the primary key that handles at least
one of those relationships.
Second, we have two uniqueness requirements for Customers:
username and email address. The username is used for actual
customer lookups, whereas the email address is solely a
requirement around uniqueness.
We’re going to focus on handling the uniquness requirements first
because that must be built into the primary key of the main table.
You can’t handle this via a secondary index.
We discussed uniqueness on two attributes in Chapter 16. You can’t
build uniqueness on multiple attributes into a single item, as that
would only ensure the combination of the attributes is unique.
Rather, we’ll need to make multiple items.
Let’s create two types of items—Customers and CustomerEmails—
with the following primary key patterns:
Customer:
• PK: CUSTOMER#<Username>
• SK: CUSTOMER#<Username>
CustomerEmail:
307
• PK: CUSTOMEREMAIL#<Email>
• SK: CUSTOMEREMAIL#<Email>
We can load our table with some items that look like the following:
So far, our service has two customers: Alex DeBrie and Vito
Corleone. For each customer, there are two items in DynamoDB.
One item tracks the customer by username and includes all
information about the customer. The other item tracks the
customer by email address and includes just a few attributes to
identify to whom the email belongs.
While this table shows the CustomerEmail items, I will hide them
when showing the table in subsequent views. They’re not critical to
the rest of the table design, so hiding them will de-clutter the table.
We can update our entity chart to add the CustomerEmails item
type and to fill out the primary key patterns for our first two items:
Entity PK SK
Customers CUSTOMER#<Username> CUSTOMER#<Username>
CustomerEmails CUSTOMEREMAIL#<Email> CUSTOMEREMAIL#<Email>
Addresses
Orders
OrderItems
308
Table 17. E-commerce entity chart
Finally, when creating a new customer, we’ll want to only create the
customer if there is not an existing customer with the same
username and if this email address has not been used for another
customer. We can handle that using a DynamoDB Transaction.
The code below shows the code to create a customer with proper
validation:
response = client.transact_write_items(
  TransactItems=[
  {
  'Put': {
  'TableName': 'EcommerceTable',
  'Item': {
  'PK': { 'S': 'CUSTOMER#alexdebrie' },
  'SK': { 'S': 'CUSTOMER#alexdebrie' },
  'Username': { 'S': 'alexdebrie' },
  'Name': { 'S': 'Alex DeBrie' },
  ... other attributes ...
  },
  'ConditionExpression': 'attribute_not_exists(PK)
  }
  },
  {
  'Put': {
  'TableName': 'EcommerceTable',
  'Item': {
  'PK': { 'S': 'CUSTOMEREMAIL#alexdebrie1@gmail.com' },
  'SK': { 'S': 'CUSTOMEREMAIL#alexdebrie1@gmail.com' },
  },
  'ConditionExpression': 'attribute_not_exists(PK)
  }
  }
  ]
)
Our TransactWriteItems API has two write requests: one to write
the Customer item and one to write the CustomerEmail item.
Notice that both have condition expressions to confirm that there is
not an existing item with the same PK. If one of the conditions is
violated, it means that either the username or email address is
already in use and thus the entire transaction will be cancelled.
Now that we’ve handled our Customer item, let’s move on to one of
309
our relationships. I’ll go with Addresses next.
19.3.2. Modeling the Addresses entity
There is a one-to-many relationship between Customers and
Addresses. We can use strategies from Chapter 11 to see how we can
handle the relationship.
The first thing we should ask is whether we can denormalize the
relationship. When denormalizing by using a complex attribute, we
need to ask two things:
1. Do we have any access patterns that fetch related entity directly
by values of the related entity, outside the context of the parent?
2. Is the amount of data in the complex attribute unbounded?
In this case, the answer to the first question is 'No'. We will show
customers their saved addresses, but it’s always in the context of the
customer’s account, whether on the Addresses page or the Order
Checkout page. We don’t have an access pattern like "Fetch
Customer by Address".
The answer to the second question is (or can be) 'No' as well. While
we may not have considered this limitation upfront, it won’t be a
burden on our customers to limit them to only 20 addresses. Notice
that data modeling can be a bit of a dance. You may not have
thought to limit the number of saved addresses during the initial
requirements design, but it’s easy to add on to make the data
modeling easier.
Because both answers were 'No', we can use the denormalization
strategy and use a complex attribute. Let’s store each customer’s
addresses on the Customer item.
Our updated table looks as follows:
310
Notice that our Customer items from before have an Addresses
attribute outlined in red. The Addresses attribute is of the map
type and includes one or more named addresses for the customer.
We can update our entity chart as follows:
Entity PK SK
Customers CUSTOMER#<Username> CUSTOMER#<Username>
CustomerEmails CUSTOMEREMAIL#<Email> CUSTOMEREMAIL#<Email>
Addresses N/A N/A
Orders
OrderItems
Table 18. E-commerce entity chart
There is no separate Address item type, so we don’t have a PK or SK
pattern for that entity.
19.3.3. Modeling Orders
Now let’s move on to the Order item. This is our second one-tomany relationship with the Customer item.
Let’s walk through the same analysis as we did with Addresses—can
we handle this relationship through denormalization?
Unfortunately, it doesn’t appear to be a good idea here. Because
DynamoDB item sizes are limited to 400KB, you can only
311
denormalize and store as a complex attribute if there is a limit to
the number of related items. However, we don’t want to limit the
number of orders that a customer can make with us—we would be
leaving money on the table! Because of that, we’ll have to find a
different strategy.
Notice that we have a join-like access pattern where we need to
fetch both the Customer and the Orders in a single request. The
next strategy, and the most common one for one-to-many
relationships, is to use the primary key plus the Query API to 'prejoin' our data.
The Query API can only fetch items with the same partition key, so
we need to make sure our Order items have the same partition key
as the Customer items. Further, we want to retrieve our Orders by
the time they were placed, starting with the most recent.
Let’s use the following pattern for our Order items:
• PK: CUSTOMER#<Username>
• SK: #ORDER#<OrderId>
For the OrderId, we’ll used a KSUID. KSUIDs are unique identifiers
that include a timestamp in the beginning. This allows for
chronological ordering as well as uniqueness. You can read more
about KSUIDs in Chapter 14.
We can add a few Order items to our table to get the following:
312
We’ve added three Order items to our table, two for Alex DeBrie
and one for Vito Corleone. Notice that the Orders have the same
partition key and are thus in the same item collection as the
Customer. This means we can fetch both the Customer and the
most recent Orders in a single request.
An additional note—see that we added a prefix of # to our Order
items. Because we want the most recent Orders, we will be fetching
our Orders in descending order. This means our Customer item
needs to be after all the Order items so that we can fetch the
Customer item plus the end of the Order items. If we didn’t have
the # prefix for Order items, then Orders would show up after the
Customer and would mess up our ordering.
To handle our pattern to retrieve the Customer and the most recent
Orders, we can write the following Query:
resp = client.query(
  TableName='EcommerceTable',
  KeyConditionExpression='#pk = :pk',
  ExpressionAttributeNames={
  '#pk': 'PK'
  },
  ExpressionAttributeValues={
  ':pk': { 'S': 'CUSTOMER#alexdebrie' }
  },
  ScanIndexForward=False,
  Limit=11
)
313
We use a key expression that uses the proper PK to find the item
collection we want. Then we set ScanIndexForward=False so that it
will start at the end of our item collection and go in descending
order, which will return the Customer and the most recent Orders.
Finally, we set a limit of 11 so that we get the Customer item plus
the ten most recent orders.
We can update our entity chart as follows:
Entity PK SK
Customers CUSTOMER#<Username> CUSTOMER#<Username>
CustomerEmails CUSTOMEREMAIL#<Email> CUSTOMEREMAIL#<Email>
Addresses N/A N/A
Orders CUSTOMER#<Username> #ORDER#<OrderId>
OrderItems
Table 19. E-commerce entity chart
19.3.4. Modeling the Order Items
The final entity we need to handle is the OrderItem. An OrderItem
refers to one of the items that was in an order, such as a specific
book or t-shirt.
There is a one-to-many relationship between Orders and
OrderItems, and we have an access pattern where we want join-like
functionality as we want to fetch both the Order and all its
OrderItems for the OrderDetails page.
Like in the last pattern, we can’t denormalize these onto the Order
item as the number of items in an order is unbounded. We don’t
want to limit the number of items a customer can include in an
order.
Further, we can’t use the same strategy of a primary key plus Query
API to handle this one-to-many relationship. If we did that, our
314
OrderItems would be placed between Orders in the base table’s
item collections. This would significantly reduce the efficiency of
the "Fetch Customer and Most Recent Orders" access pattern we
handled in the last section as we would now be pulling back a ton of
extraneous OrderItems with our request.
That said, the principles we used in the last section are still valid.
We’ll just handle it in a secondary index.
First, let’s create the OrderItem entity in our base table. We’ll use
the following pattern for OrderItems:
• PK: ORDER#<OrderId>#ITEM#<ItemId>
• SK: ORDER#<OrderId>#ITEM#<ItemId>
Our table will look like this:
We have added two OrderItems into our table. They are outlined in
red at the bottom. Notice that the OrderItems have the same
OrderId as an Order in our table but that the Order and OrderItems
are in different item collections.
To get them in the same item collection, we’ll add some additional
properties to both Order and OrderItems.
315
The GSI1 structure for Orders will be as follows:
• GSI1PK: ORDER#<OrderId>
• GSI1SK: ORDER#<OrderId>
The GSI1 structure for OrderItems will be as follows:
• GSI1PK: ORDER#<OrderId>
• GSI1SK: ITEM#<ItemId>
Now our base table looks as follows:
Notice that our Order and OrderItems items have been decorated
with the GSI1PK and GSI1SK attributes.
We can then look at our GSI1 secondary index:
316
Now our Orders and OrderItems have been re-arranged so they are
in the same item collection. As such, we can fetch an Order and all
of its OrderItems by using the Query API against our secondary
index.
The code to fetch an Order and all of its OrderItems is as follows:
resp = client.query(
  TableName='EcommerceTable',
  IndexName='GSI1',
  KeyConditionExpression='#gsi1pk = :gsi1pk',
  ExpressionAttributeNames={
  '#gsi1pk': 'GSI1PK'
  },
  ExpressionAttributeValues={
  ':gsi1pk': 'ORDER#1VrgXBQ0VCshuQUnh1HrDIHQNwY'
  }
)
We can also update our entity chart as follows:
Entity PK SK
Customers CUSTOMER#<Username> CUSTOMER#<Username>
CustomerEmails CUSTOMEREMAIL#<Email> CUSTOMEREMAIL#<Email>
317
Entity PK SK
Addresses N/A N/A
Orders CUSTOMER#<Username> #ORDER#<OrderId>
OrderItems ORDER#<OrderId>#ITEM#
<ItemId>
ORDER#<OrderId>#ITEM#
<ItemId>
Table 20. E-commerce entity chart
Further, let’s make a corresponding entity chart for GSI1 so we can
track items in that index:
Entity GSI1PK GSI1SK
Customers
CustomerEmails
Addresses
Orders ORDER#<OrderId> ORDER#<OrderId>
OrderItems ORDER#<OrderId> ITEM#<ItemId>
Table 21. E-commerce GSI1 entity chart
19.4. Conclusion
We’re starting to get warmed up with our DynamoDB table design.
In this chapter, we looked at some advanced patterns including
using primary key overloading to create item collections with
heterogeneous items.
Let’s review our final solution.
Table Structure
Our table uses a composite primary key with generic names of PK
and SK for the partition key and sort key, respectively. We also have
a global secondary index named GSI1 with similarly generic names
of GSI1PK and GSI1SK for the partition and sort keys.
318
Our final entity chart for the main table is as follows:
Entity PK SK
Customers CUSTOMER#<Username> CUSTOMER#<Username>
CustomerEmails CUSTOMEREMAIL#<Email> CUSTOMEREMAIL#<Email>
Addresses N/A N/A
Orders CUSTOMER#<Username> #ORDER#<OrderId>
OrderItems ORDER#<OrderId>#ITEM#
<ItemId>
ORDER#<OrderId>#ITEM#
<ItemId>
Table 22. E-commerce entity chart
And the final entity chart for the GSI1 index is as follows:
Entity GSI1PK GSI1SK
Customers
CustomerEmails
Addresses
Orders ORDER#<OrderId> ORDER#<OrderId>
OrderItems ORDER#<OrderId> ITEM#<ItemId>
Table 23. E-commerce GSI1 entity chart
Notice a few divergences from our ERD to our entity chart. First,
we needed to add a special item type, 'CustomerEmails', that are
used solely for tracking the uniqueness of email addresses provided
by customers. Second, we don’t have a separate item for Addresses
as we denormalized it onto the Customer item.
After you make these entity charts, you should include them in the
documentation for your repository to assist others in knowing how
the table is configured. You don’t want to make them dig through
your data access layer to figure this stuff out.
Access Patterns
We have the following six access patterns that we’re solving:
319
Access Pattern Index Parameters Notes
Create Customer N/A N/A
Use TransactWriteItems to
create Customer and
CustomerEmail item with
conditions to ensure
uniqueness on each
Create / Update
Address N/A N/A
Use UpdateItem to update
the Addresses attribute on
the Customer item
View Customer &
Most Recent
Orders
Main table • Username
Use
ScanIndexForward=False
to fetch in descending
order.
Save Order N/A N/A
Use TransactWriteItems to
create Order and
OrderItems in one request
Update Order N/A N/A Use UpdateItem to update
the status of an Order
View Order &
Order Items GSI1
• OrderId
Table 24. E-commerce access patterns
Just like your entity charts, this chart with your access pattern
should be included in the documentation for your repository so
that it’s easier to understand what’s happening in your application.
In the next chapter, we’re going to dial it up to 11 by modeling a
complex application with a large number of entities and
relationships.

