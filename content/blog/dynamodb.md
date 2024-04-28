+++
title = "DynamoDB"
date = "2024-04-12"
tags = ["AWS"]
subtitle = "Some notes from the The DynamoDB Book"
+++

<!--
Done: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13, 14, 16, 17, 18, 19
In-progress: 20
Omitted: 12, 15, 22
Remaining: 21
-->

## Overview
- DynamoDB is a fully-managed NoSQL database.
- Two data models are supported:
  - *Key-value*: distributed hash table; store and retrieve records by primary key.
  - *Wide-column*: hash table where the value of each record is a B-tree; allows for performing contiguous range queries.
- It provides fast, consistent performance that does not degrade as the amount of stored data increases.
- Operations are performed by making HTTP requests which supports a virtually unlimited number of concurrent requests.
- A flexible pricing model is available:
  - *Provisioned*: read/write throughput configured independently and dynamically.
  - *On-demand*: pay per request.
- DynamoDB is suited for most OLTP applications with high volumes of small records being read and written.
- The rigidity of operations prevents writing inefficient queries and producing data models that do not scale.
- Filtering is built into the data model as the primary keys and secondary indexes determine how data is retrieved.

## Concepts

### Tables, Items & Attributes
- A table is a collection of items, each composed of one or more attributes.
- An attribute is given a type when it is written:
  - *Scalar*: single simple value (*string*, *number*, *binary*, *boolean*, *null*).
  - *Complex*: groupings with arbitrary nested attributes (*list*, *map*).
  - *Set*: multiple unique values of the same type (*string*, *number*, *binary* sets).
- The attribute type affects which operations can be performed.
- A table is schemaless and attributes are not required on every item.
- Multiple different types of entities are commonly stored in a single table to allow for handling complex access patterns in a single request without the need for joins.
  - Use prefixes to differentiate between different entity types (e.g. `CUSTOMER#123`).

![items.png](/img/dynamodb/items.png)

### Primary Keys
- A primary key must be declared when creating a table.
- Each item must include the primary key and is uniquely identified by that key.
- There are two types of primary keys:
  - *Simple*: consists of a single partition key (hash).
  - *Composite*: consists of a partition key and sort key (range).
- A simple primary key allows for only fetching a single item at a time while a composite primary key enables fetching all items with the same partition key.
- Writing an item using a primary key that already exits will overwrite the existing item (unless this is explictly disabled, in which case the write will be rejected).

### Secondary Indexes
- Secondary indexes allow for reshaping data into another format for querying.
- All items will be copied from the base table into the secondary index in the reshaped form.
- A key schema consisting of a partition key and (optional) sort key must be declared when creating a secondary index.
- There are two types of secondary indexes:
  - *Local*: uses the same partition key as the primary key but a different sort key.
  - *Global*: uses any attributes for the partition key and sort key.
- Separate application attributes (*FirstName*, *LastUdated*, etc.) from indexing attributes (*PK*, *GSI1SK*, etc.) and avoid reusing attribute across multiple indexes to avoid unnecessary complexity.

| | LSI | GSI |
|-|-|-|
| Throughput | Shared with base table | Separately provisioned |
| Consistency | Allow opting into strongly-consistent reads | Eventually-consistent reads only |
| Creation time | Must be specified when table created | Created and deleted as needed |

### Item Collections
- An item collection consists of the group of items that share the same partition key in either the base table or a secondary index.
- All the items in an item collection are ordered and stored as a B-tree and will be allocated to the same partition.
- The item collection size includes both the size of the items in the base table and the size of the items in the local secondary index.

![item-collections.png](/img/dynamodb/item-collections.png)

### Attribute Projections
- A projection is the set of attributes that is copied into a secondary index.
- The partition key and sort key of the table are always projected into the index and other attributes can be projected as required.
- This must be declared when the secondary index is created:
  - *KEYS_ONLY*: project the table partition key and sort key values, plus the index key values.
  - *INCLUDE*: include other additional non-key attributes.
  - *ALL*: include all the attributes from the base table.
- Choosing which attributes to project presents a trade-off between throughput costs and storage costs.
- It is not possible to change projected attributes once the index has been created.

### Sparse Indexes
- When an item is written to the base table, it will be copied into the secondary index only if it has the elements of the key schema for the secondary index.
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

### Sharding
- Sharding refers to splitting data across multiple partitions to prevent hot partitions.
- There are a number of different read-sharding patterns:
  - Random partitions with scatter-gather at read time.
  - Grouping items into item collections based on similar time range.

### Consistency
- Consistency refers to whether a read operation receives all prior write operations.
- There are two consistency options available:
  - *Strong*: any item read will reflect all prior writes.
  - *Eventual*: it is possible that read items will not reflect all prior writes.
- Reads are eventually-consistent by default though it is possible to opt into strongly-consistent for base tables and LSIs.
- Sources of eventual consistency include asynchronous data replication from primary to secondary nodes and from base tables to GSIs.

### Key Overloading 
- Key overloading refers to using generic names for the primary key attributes and using different values dependending on the item type.
- Prefixes are used for the partition and sort key values in order to identify the item type and to avoid overlap between different item types.
- Secondary indexes can be overloaded just like primary keys.
- Common generic names are: *PK*, *SK*, *GSI1PK*, *GSI1SK*, etc.

![key-overloading.png](/img/dynamodb/key-overloading.png)

### Time-to-Live
- TTLs allow for automatically deleting items after a specified time.
- An attribute of type *number* is specified on the table that will serve as the marker for item deletion.
- To expire an item, a Unix timestamp at seconds granularity can be stored in the specified attribute that indicates that time after which the item should be deleted.
- Items are usually deleted within 48 hours after the time indicated by the attribute.

### Capacity Units
- A RCU represents one strongly-consistent read per second or two eventually-consistent reads per second for an item up to 4KB.
  - Transactional read requests require two RCUs to perform one read per second for items up to 4KB.
  - Reading an item larger than 4KB will consume additional RCUs.
- One WCU represents one write per second for an item up to 1KB in size.
  - Transactional write requests require two WCUs to perform one write per second for items up to 1KB.
  - Writing an item larger than 1KB will consume additional WCUs.

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

### Batch/Transact
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
  - Expression attribute values: represent an attribute value being evaluated in the request (e.g. `:foo`).
  - Expression attribute names: specify the names of attributes (e.g. `#bar`).
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

### Session Store

#### Entity Relationships

![session-store-1.png](/img/dynamodb/session-store-1.png)

#### Entity Chart
| Entity  | PK               | SK |
|---------|------------------|----|
| Session | `<SessionToken>` | -  |
| User    | -                | -  |

| Entity  | GSI1PK       | GSI1SK |
|---------|--------------|--------|
| Session | `<Username>` | -      |
| User    | -            | -      |

#### Access Patterns
| Access Pattern              | Index      | Parameters   | Notes                                          |
|-----------------------------|------------|--------------|------------------------------------------------|
| Create session              | Base table | SessionToken | Use condition expression to enforce uniqueness |
| Get session                 | Base table | SessionToken | Use filter expression to handle expired tokens |
| Delete session (time-based) | N/A        | N/A          | Handled by TTL                                 |
| Delete sessions for user    | UserIndex  | Username     | Find tokens for user                           |
|                             | Main table | SessionToken | Delete each token                              |

#### Table Structure

![session-store-2.png](/img/dynamodb/session-store-2.png)

```py
# prevent duplicate session tokens from being written using a condition expression
created_at = datetime.datetime.now()
expires_at = created_at + datetime.timedelta(days=7)
response = client.put_item(
  TableName='SessionStore',
  Item={
    'SessionToken': { 'S': str(uuid.uuidv4()) },
    'Username': { 'S': 'dave' },
    'CreatedAt': { 'S': created_at.isoformat() },
    'ExpiresAt': { 'S': expires_at.isoformat() },
  },
  ConditionExpression='attribute_not_exists(SessionToken)',
)
```
```py
# handle expired items using a filter expression
epoch_seconds = int(time.time())
results = client.query(
  TableName='SessionStore',
  KeyConditionExpression='#token = :token',
  FilterExpression='#ttl >= :epoch',
  ExpressionAttributeNames={
    '#token': 'SessionToken',
    '#ttl': 'TTL'
  },
  ExpressionAttributeValues={
    ':token': { 'S': '0bc6bdf8-6dac-4212-b11a-81f784297c78' },
    ':epoch': { 'N': str(epoch_seconds) }
  }
)
```

![session-store-3.png](/img/dynamodb/session-store-3.png)

```py
# query GSI for session tokens related to user and delete each one
results = client.query(
  TableName='SessionStore',
  IndexName='UserIndex',
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

#### Entity Chart

| Entity         | PK                              | SK                              |
|----------------|---------------------------------|---------------------------------|
| Customers      | `CUSTOMER#<Username>`           | `CUSTOMER#<Username>`           |
| CustomerEmails | `CUSTOMEREMAIL#<Email>`         | `CUSTOMEREMAIL#<Email>`         |
| Addresses      |                                 |                                 |
| Orders         | `CUSTOMER#<Username>`           | `#ORDER#<OrderId>`              |
| OrderItems     | `ORDER#<OrderId>#ITEM#<ItemId>` | `ORDER#<OrderId>#ITEM#<ItemId>` |

| Entity         | GSI1PK            | GSI1SK            |
|----------------|-------------------|-------------------|
| Customers      |                   |                   |
| CustomerEmails |                   |                   |
| Addresses      |                   |                   |
| Orders         | `ORDER#<OrderId>` | `ORDER#<OrderId>` |
| OrderItems     | `ORDER#<OrderId>` | `ITEM#<ItemId>`   |

#### Table Structure

![e-commerce-2.png](/img/dynamodb/e-commerce-2.png)

![e-commerce-3.png](/img/dynamodb/e-commerce-3.png)

![e-commerce-4.png](/img/dynamodb/e-commerce-4.png)

![e-commerce-5.png](/img/dynamodb/e-commerce-5.png)

![e-commerce-6.png](/img/dynamodb/e-commerce-6.png)

#### Access Patterns

| Access pattern | Index | Parameters | Notes |
|-|-|-|-|
| Create customer               | N/A        | N/A      | Use *TransactWriteItems* to create *Customer* and *CustomerEmail* to ensure uniqueness |
| Create/update addresss        | N/A        | N/A      | Use *UpdateItem* to update the *Addresses* attribute on the Customer item              |
| View customer & recent orders | Base table | Username | Use `ScanIndexForward=False` to fetch in descending order                              |
| Save order                    | N/A        | N/A      | Use *TransactWriteItems* to create *Order* and *OrderItems*                            |
| Update order                  | N/A        | N/A      | Use *UpdateItem* to update the status of an *Order*                                    |
| View order & order items      | GSI1       | OrderId  |                                                                                        |

#### Code Snippets
```py
# ensure uniqueness of both the username and email address
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
          # ...
        },
        'ConditionExpression': 'attribute_not_exists(PK)',
      },
    },
    {
      'Put': {
        'TableName': 'EcommerceTable',
        'Item': {
          'PK': { 'S': 'CUSTOMEREMAIL#alexdebrie1@gmail.com' },
          'SK': { 'S': 'CUSTOMEREMAIL#alexdebrie1@gmail.com' },
        },
        'ConditionExpression': 'attribute_not_exists(PK)',
      },
    },
  ]
)

# retrieve customer and most recent orders
response = client.query(
  TableName='EcommerceTable',
  KeyConditionExpression='#pk = :pk',
  ExpressionAttributeNames={
    '#pk': 'PK'
  },
  ExpressionAttributeValues={
    ':pk': { 'S': 'CUSTOMER#alexdebrie' }
  },
  # start at the end of the item collection and read in descending order
  ScanIndexForward=False,
  Limit=11
)

# fetch an order and all its items
response = client.query(
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
```


### Big Time Deals

#### Entity Relationships

![big-time-deals-1.png](/img/dynamodb/big-time-deals-1.png)

#### Entity Chart
| Entity       | PK              | SK              |
|-|-|-|
| Deal         | `DEAL#<DealId>` | `DEAL#<DealId>` |
| Brand        | `BRAND#<Brand>` | `BRAND#<Brand>` |
| Brand Like   | `BRANDLIKE#<Brand>#<Username>` | `BRANDLIKE#<Brand>#<Username>` |
| Category     |    | |
| FeaturedDeal |    | |
| Page         |    | |
| User         |    | |
| Message      |    | |


| Entity       | GSI1PK                       | GSI1SK          |
|-|-|-|
| Deal         | `DEALS#<TruncatedTimestamp>` | `DEAL#<DealId>` |
| Brand        |    | |
| Brand Like   |    | |
| Category     |    | |
| FeaturedDeal |    | |
| Page         |    | |
| User         |    | |
| Message      |    | |

| Entity       | GSI2PK                               | GSI1SK          |
|-|-|-|
| Deal         | `BRAND#<Brand>#<TruncatedTimestamp>` | `DEAL#<DealId>` |
| Brand        |    | |
| Category     |    | |
| FeaturedDeal |    | |
| Page         |    | |
| User         |    | |
| Message      |    | |

| Entity       | GSI2PK                                     | GSI1SK          |
|-|-|-|
| Deal         | `CATEGORY#<Category>#<TruncatedTimestamp>` | `DEAL#<DealId>` |
| Brand        |    | |
| Category     |    | |
| FeaturedDeal |    | |
| Page         |    | |
| User         |    | |
| Message      |    | |



#### Table Structure

##### Deal

![big-time-deals-2.png](/img/dynamodb/big-time-deals-2.png)

![big-time-deals-3.png](/img/dynamodb/big-time-deals-3.png)

```py
def fetch_items_for_date(date, last_seen='$', limit=25):
  response = client.query(
    TableName='BigTimeDeals',
    Index='GSI1',
    KeyConditionExpression='#pk = :pk AND #sk < :sk',
    ExpressionAttributeNames={
      '#pk': 'PK',
      '#sk': 'SK'
    },
    ExpressionAttributeValues={
      ':pk': { 'S': f"DEALS#{date.strftime('%Y-%m-%dT00:00:00')}" },
      'sk': { 'S': f'DEAL#{last_seen}' }
    },
    ScanIndexForward=False,
    Limit=limit
  )

# fetch latest 25 deals
def get_deals(date, last_deal_seen='$'):
  deals = []
  count = 0
  while len(deals) < 25 and count < 5:
    items = fetch_items_for_date(date, last_deal_seen)
    for item in items:
      deals.append(Deal(
        title=item['Title']['S'],
        deal_id=item['DealId']['S'],
        link=item['Link']['S'],
        # ...
      ))
    date = date - datetime.timedelta(days=1)
    count += 1

  return deals[:24]
```

![big-time-deals-4.png](/img/dynamodb/big-time-deals-4.png)

```py
# retrieve cached item from random partition
shard = random.randint(1, 10)
response = client.get_item(
  TableName='BigTimeDeals',
  Key={
    'PK': { 'S': f"DEALSCACHE#{shard}" },
    'SK': { 'S': f"DEALSCACHE#{shard}" },
  }
)
```

##### Brand

![big-time-deals-5.png](/img/dynamodb/big-time-deals-5.png)

```py
response = dynamodb.transact_write_items(
  TransactItems=[
    {
      'Put': {
        'TableName': 'BigTimeDeals',
        'Item': {
          'PK': { 'S': 'BRANDLIKE#APPLE#testuser' },
          'SK': { 'S': 'BRANDLIKE#APPLE#testuser' },
          # ...
        },
        'ConditionExpression': 'attribute_not_exists(PK)',
      }
    },
  ]
)
```

  {
  "Update": {
  "Key": {
  "PK": { "S": "BRAND#APPLE" },
  "SK": { "S": "BRAND#APPLE" },
  },
  "TableName": "BigTimeDeals",
  "ConditionExpression": "attribute_exists(PK)"
  "UpdateExpression": "SET #likes = #likes + :incr",
  "ExpressionAttributeNames": {
  "#likes": "LikesCount"
  },
  "ExpressionAttributeValues": {
  ":incr": { "N": "1" }
  }
  }
  }
  ]


The first tries to create a Brand Like item for the User "alexdebrie" for the Brand "Apple". Notice that it includes a condition expression to ensure the Brand Like item doesn’t already exist, which would indicate the user already liked this Brand.  The second operation increments the LikesCount on the Apple Brand item by 1. It also includes a condition expression to ensure the Brand exists.  Remember that each operation will succeed only if the other operation also succeeds. If the Brand Like already exists, the LikesCount won’t be incremented. If the Brand doesn’t exist, the Brand Like won’t be created.

Handling Brand Watches
In addition to liking a Brand, users may also watch a Brand. When a
user is watching a Brand, they will also be notified whenever a new
Deal is added for a Brand.
The patterns here are pretty similar to the Brand Like discussion,
with one caveat: our internal system needs to be able to find all
watchers for a Brand so that we can notify them.
Because of that, we’ll put all Brand Watch items in the same item
collection so that we can run a Query operation on it.
We’ll use the following primary key pattern:
• PK: BRANDWATCH#<Brand>
• SK: USER#<Username>
Notice that Username is not a part of the partition key in this one,
so all Brand Watch items will be in the same item collection.
We’ll still use the same DynamoDB Transaction workflow when
adding a watcher to increase the WatchCount while ensuring that
the user has not previously watched the Brand.
Sending New Brand Deal Messages to Brand Watchers
350
One of our access patterns is to send a New Brand Deal Message to
all Users that are watching a Brand. Let’s discuss how to handle that
here.
Remember that we can use DynamoDB Streams to react to changes
in our DynamoDB table. To handle this access pattern, we will
subscribe to the DynamoDB stream. Whenever we receive an event
on our stream, we will do the following:
1. Check to see if the event is indicating that a Deal item was
inserted into our table;
2. If yes, use the Query operation to find all BrandWatch items for
the Brand of the new Deal;
3. For each BrandWatch item found, create a new Message item for
the user that alerts them of the new Brand Deal.
For steps 2 & 3, the code will be similar to the following:
resp = dynamodb.query(
  TableName='BigTimeDeals',
  KeyConditionExpression="#pk = :pk",
  ExpressionAttributeNames={
  '#pk': "PK"
  },
  ExpressionAttributeValues={
  ':pk': { 'S': "BRANDWATCH#APPLE" }
  }
)
for item in resp['Items']:
  username = item['Username']['S']
  send_message_to_brand_watcher(username)
First, we run a Query operation to fetch all the watchers for the
given Brand. Then we send a message to each watcher to alert them
of the new deal.
Patterns like this make it much easier to handle reactive
functionality without slowing down the hot paths of your
application.
351
20.3.3. Modeling the Category item
Now that we’ve handled the Brand item, let’s handle the Category
item as well. Categories are very similar to Brands, so we won’t
cover it in quite as much detail.
There are two main differences between the Categories and Brands:
1. Categories do not have a "Fetch all Categories" access pattern, as
there are only eight categories that do not change.
2. The Categories page has a list of 5-10 "Featured Deals" within the
Category.
For the first difference, that means we won’t have a singleton
"CATEGORIES" item like we had with Brands.
For the second difference, we need to figure out how to indicate a
Deal is 'Featured' within a particular Category. One option would be
to add some additional attributes and use a sparse index pattern to
group a Category with its Featured Deals. However, that seems a bit
overweight. We know that we only have a limited number of
Featured Deals within a particular Category.
Instead, let’s combine our two denormalization strategies from the
one-to-many relationships chapter. We’ll store information about
Featured Deals in a complex attribute on the Category itself, and
this information will be duplicated from the underlying Deal item.
Our Category items would look as follows:
352
Each Category item includes information about the Featured Deals
on the item directly. Remember that setting Featured Deals is an
internal use case, so we can program our internal CMS such that it
includes all information about all Featured Deals whenever an
editor is setting the Featured Deals for a Category.
Thus, for modeling out the Category items and related entities, we
create the following item types:
Category
• PK: CATEGORY#<Category>
• SK: CATEGORY#<Category>
CategoryLike
• PK: CATEGORYLIKE#<Category>#<Username>
• SK: CATEGORYLIKE#<Category>#<Username>
CategoryWatch
• PK: CATEGORYWATCH#<Category>
• SK: USER#<Username>
20.3.4. Featured Deals and Editors' Choice
We’re almost done with the Deals-related portion of this data
353
model. The last thing we need to handle is around the Featured
Deals on the front page of the application and the Editor’s Choice
page.
For me, this problem is similar to the "Featured Deals for Category"
problem that we just addressed. We could add some attributes to a
Deal to indicate that it’s featured and shuttle them off into a
secondary index. In contrast, we could go a much simpler route by
just duplicating some of that data elsewhere.
Let’s use a combination of the 'singleton item' strategy we discussed
earlier with this duplication strategy. We’ll create two new singleton
items: one for the Front Page and one for the Editor’s Choice page.
The table might look as follows:
Notice the singleton items for the Front Page and the Editor’s
Choice page. Additionally, like we did with the Deals Cache items,
we could copy those across a number of partitions if needed. This is
a simple, effective way to handle these groupings of featured deals.
Let’s take a breath here and take a look at our updated entity chart
with all Deal-related items finished.
354
Entity PK SK
Deal DEAL#<DealId> DEAL#<DealId>
Brand BRAND#<Brand> BRAND#<Brand>
Brands BRANDS BRANDS
BrandLike BRANDLIKE#<Brand>#<Us
ername>
BRANDLIKE#<Brand>#<Us
ername>
BrandWatch BRANDWATCH#<Brand> USER#<Username>
Category CATEGORY#<Category> CATEGORY#<Category>
CategoryLike CATEGORYLIKE#<Categor
y>#<Username>
CATEGORYLIKE#<Categor
y>#<Username>
CategoryWatch CATEGORYWATCH#<Catego
ry>
USER#<Username>
FrontPage FRONTPAGE FRONTPAGE
Editor’s Choice EDITORSCHOICE EDITORSCHOICE
User
Message
Table 28. Big Time Deals entity chart
We also have some attributes in our secondary indexes. To save
space, I’ll only show items that have attributes in those indexes.
First, our GSI1 secondary index:
Entity GSI1PK GSI1SK
Deal DEALS#<TruncatedTimestamp> DEAL#<DealId>
Table 29. Big Time Deals GSI1 entity chart
Then, our GSI2 secondary index:
Entity GSI2PK GSI2SK
Deal BRAND#<Brand>#<TruncatedTi
mestamp>
DEAL#<DealId>
Table 30. Big Time Deals GSI2 entity chart
And finally, the GSI3 secondary index:
355
Entity GSI3PK GSI3SK
Deal CATEGORY#<Category>#<Trunc
atedTimestamp>
DEAL#<DealId>
Table 31. Big Time Deals GSI3 entity chart
20.3.5. Modeling the User item
Now that we have most of the entities around Deals modeled, let’s
move on to Users and Messages.
In addition to the simple Create / Read / Update User access
patterns, we have the following access patterns that are based on
Users:
• Fetch all Messages for User
• Fetch all Unread Messages for User
• Mark Message as Read
• Send new Hot Deal Message to all Users
Remember that we handled the "Send new Brand Deal Message to
all Brand Watchers" and "Send new Category Deal Message to all
Category Watchers" in the sections on Brands and Categories,
respectively.
I’m going to start with the last access pattern—send new Hot Deal
Message to All Users—then focus on the access patterns around
fetching Messages.
The User item and Finding all Users
When we think about the "Send new Hot Deal Message to all Users",
it’s really a two-step operation:
356
1. Find all Users in our application
2. For each User, send a Message
For the 'find all Users' portion of it, we might think to mimic what
we did for Brands: use a singleton item to hold all usernames in an
attribute. However, the number of Users we’ll have is unbounded. If
we want our application to be successful, we want to have as many
Users as possible, which means we want to exceed 400KB of data.
A second approach we could do is to put all Users into a single
partition. For example, we could have a secondary index where
each User item had a static partition key like USERS so they were all
grouped together. However, this could lead to hot key issues. Each
change to a User item would result in a write to the same partition
in the secondary index, which would result in a huge number of
writes.
Instead, let’s use one of our sparse indexing strategies from Chapter
13. Here, we want to use the second type of sparse index, which
projects only a single type of entity into a table.
To do this, we’ll create a User entity with the following attributes:
• PK: USER#<Username>
• SK: USER#<Username>
• UserIndex: USER#<Username>
Our table with some User items will look as follows:
357
Notice that we have three User items in our table. I’ve also placed a
Deal item to help demonstrate how our sparse index works.
Each of our User items has a UserIndex attribute. We will create a
secondary index that uses the UserIndex attribute as the partition
key. That index looks as follows:
Notice that only our User items have been copied into this table.
Because other items won’t have that attribute, this is a sparse index
containing just Users.
Now if we want to message every User in our application, we can
use the following code:
resp = dynamodb.scan(
  TableName='BigTimeDeals',
  IndexName='UserIndex'
)
for item in resp['Items']:
  username = item['Username']['S']
  send_message_to_user(username)
358
This is a similar pattern to what we did when sending messages to
Brand or Category Watchers. However, rather than using the Query
operation on a partition key, we’re using the Scan operation on an
index. We will scan our index, then message each User that we find
in our index.
One final note on sparse indexes: an astute observer might note
that all of our secondary indexes are sparse indexes. After all, only
the Deal item has been projected into GSI1, GSI2, and GSI3. What’s
the difference here?
With the UserIndex, we’re intentionally using the sparseness of the
index to aid in our filtering strategy. As we move on, we may
include other items in our GSI1 index. However, we can’t add other
items into our UserIndex as that would defeat the purpose of that
index.
Handling our User Messages
The final access patterns we need to handle are with the Messages
for a particular User. There are three access patterns here:
• Find all Messages for User
• Find all Unread Messages for User
• Mark Message as Read
Notice the first two are the same pattern with an additional filter
condition. Let’s think through our different filtering strategies from
Chapter 13 on how we handle the two different patterns.
We probably won’t want to distinguish between these using a
partition key (e.g. putting read and unread messages for a User in
different partitions). This would require twice as many reads on the
359
'Find all Messages' access pattern, and it will be difficult to find the
exact number of Messages we want without overfetching.
Likewise, we can’t use the sort key to filter here, whether we use it
directly or as part of a composite sort key. The composite sort key
works best when you always want to filter on a particular value.
Here, we sometimes want to filter and sometimes don’t.
That leaves us to two types of strategies: using a sparse index, or
overfetching and then filtering apart from the core access of
DynamoDB, either with a filter expression or with client-side
filtering. I like to avoid the overfetching strategies unless there are a
wide variety of filtering patterns we need to support. Since we only
need to provide one here, let’s go with a sparse index.
When modeling Users, we used the sparse index to project only a
particular type of entity into an index. Here, we’re going to use the
other sparse index strategy where we filter within a particular entity
type.
Let’s create our Message item with the following pattern:
• PK: MESSAGES#<Username>
• SK: MESSAGE#<MessageId>
• GSI1PK: MESSAGES#<Username>
• GSI1SK: MESSAGE#<MessageId>
For the MessageId, we’ll stick with the KSUID that we used for
Deals and discussed in Chapter 14.
Note that the PK & SK patterns are the exact same as the GSI1PK
and GSI1SK patterns. The distinction is that the GSI1 attributes will
only be added for unread Messages. Thus, GSI1 will be a sparse index
for unread Messages.
Our table with some Message items looks as follows:
360
We have four Messages in our table. They’re grouped according to
Username, which makes it easy to retrieve all Messages for a User.
Also notice that three of the Messages are unread. For those three
Messages, they have GSI1PK and GSI1SK values.
When we look at our GSI1 secondary index, we’ll see only unread
Messages for a User:
361
This lets us quickly retrieve unread Messages for a User.
The modeling part is important, but I don’t want to leave out how
we implement this in code either. Let’s walk through a few code
snippets.
First, when creating a new Message, it will be marked as Unread.
Our create_message function will handle adding the GSI1
attributes:
def create_message(message):
  resp = client.put_item(
  TableName='BigTimeDeals',
  Item={
  'PK': { 'S': f"MESSAGE#{message.username}" },
  'SK': { 'S': f"MESSAGE#{message.created_at}" },
  'Subject': { 'S': f"MESSAGE#{message.subject}" },
  'Unread': { 'S': "True" },
  'GSI1PK': { 'S': f"MESSAGE#{message.username}" },
  'GSI1SK': { 'S': f"MESSAGE#{message.created_at}" },
  }
  )
  return message
Notice that the caller of our function doesn’t need to add the GSI1
values or even think about whether the message is unread. Because
it’s unread by virtue of being new, we can set that property and
both of the GSI1 properties in the data access layer.
Second, let’s see how we would update a Message to mark it read:
362
def mark_message_read(message):
  resp = client.update_item(
  TableName='BigTimeDeals',
  Key={
  'PK': { 'S': f"MESSAGE#{message.username}" },
  'SK': { 'S': f"MESSAGE#{message.created_at}" },
  },
  UpdateExpression="SET #unread = :false, REMOVE #gsi1pk, gsi1sk",
  ExpressionAttributeNames={
  '#unread': 'Unread',
  '#gsi1pk': 'GSI1PK',
  '#gsi1sk': 'GSI1SK'
  },
  ExpressionAttributeValues={
  ':false': { 'S': 'False' }
  }
  )
  return message
In this method, we would run an UpdateItem operation to do two
things:
1. Change the Unread property to "False", and
2. Remove the GSI1PK and GSI1SK attributes so that it will be
removed from the sparse index.
Again, the calling portion of our application doesn’t need to worry
about modifying indexing attributes on our item. That is all left in
the data access portion.
Finally, let’s see our code to retrieve all messages:
363
def get_messages_for_user(username, unread_only=False):
  args = {
  'TableName': 'BigTimeDeals',
  'KeyConditionExpression': '#pk = :pk',
  'ExpressionAttributeNames': {
  '#pk': 'PK'
  },
  'ExpressionAttributeValues': {
  ':pk': { 'S': f"MESSAGE#{username}" }
  },
  'ScanIndexForward': False
  }
  if unread_only:
  args['IndexName'] = 'GSI1'
  resp = client.query(**args)
We can use the same method for fetching all Messages and fetching
unread Messages. With our unread_only argument, a caller can
specify whether they only want unread Messages. If that’s true, we’ll
add the IndexName property to our Query operation. Otherwise,
we’ll hit our base table.
With this sparse index pattern, we’re able to efficiently handle both
access patterns around Messages.



20.4. Conclusion
Entity PK SK
Deal DEAL#<DealId> DEAL#<DealId>
Brand BRAND#<Brand> BRAND#<Brand>
Brands BRANDS BRANDS
BrandLike BRANDLIKE#<Brand>#<Us
ername>
BRANDLIKE#<Brand>#<Us
ername>
BrandWatch BRANDWATCH#<Brand> USER#<Username>
Category CATEGORY#<Category> CATEGORY#<Category>
CategoryLike CATEGORYLIKE#<Categor
y>#<Username>
CATEGORYLIKE#<Categor
y>#<Username>
CategoryWatch CATEGORYWATCH#<Catego
ry>
USER#<Username>
FrontPage FRONTPAGE FRONTPAGE
Editor’s Choice EDITORSCHOICE EDITORSCHOICE
User USER#<Username> USER#<Username>
Message MESSAGES#<Username> MESSAGE#<MessageId>
Table 32. Big Time Deals entity chart
Then, our GSI1 secondary index:
365
Entity GSI1PK GSI1SK
Deal DEALS#<TruncatedTimestamp> DEAL#<DealId>
UnreadMessages MESSAGES#<Username> MESSAGE#<MessageId>
Table 33. Big Time Deals GSI1 entity chart
Then, our GSI2 secondary index:
Entity GSI2PK GSI2SK
Deal BRAND#<Brand>#<TruncatedTi
mestamp>
DEAL#<DealId>
Table 34. Big Time Deals GSI2 entity chart
And finally, the GSI3 secondary index:
Entity GSI3PK GSI3SK
Deal CATEGORY#<Category>#<Trunc
atedTimestamp>
DEAL#<DealId>




Access Pattern Index Parameters Notes
Create Deal N/A N/A Will happen in
internal CMS
Create Brand N/A N/A Add to BRANDS
container object
Create Category N/A N/A Fixed number of
categories (8)
Set Featured
Deals for Front
Page
N/A N/A
Will happen in
internal CMS. Send up
all featured deals.
Set Featured
Deals for
Category
N/A N/A
Will happen in
internal CMS. Send up
all featured deals.
Set Featured
Deals for Editor’s
Choice Page
N/A N/A
Will happen in
internal CMS. Send up
all featured deals.
366
Access Pattern Index Parameters Notes
Fetch Front Page
& Latest Deals
Main table N/A Fetch Front Page Item
GSI1 • LastDealIdSeen
Query timestamp
partitions for up to 25
deals
Fetch Category &
Latest Deals
Main table • CategoryName Fetch Category Item
GSI3
• CategoryName
• LastDealIdSeen
Query timestamp
partitions for up to 25
deals
Fetch Editor’s
Choice Page Main table N/A Fetch Editor’s Choice
item
Fetch Latest Deals
for Brand GSI2 *BrandName
Query timestamp
partitions for up to 25
deals
Fetch all Brands Main table N/A Fetch BRANDS
container item
Fetch Deal Main table • Brand GetItem on Deal Id
Create User Main table N/A
Condition expression
to ensure uniqueness
on username
Like Brand For
User Main table N/A
Transaction to
increment Brand
LikeCount and ensure
User hasn’t liked
WatchBrand For
User Main table N/A
Transaction to
increment Brand
WatchCount and
ensure User hasn’t
watched
Like Category For
User Main table N/A
Transaction to
increment Category
LikeCount and ensure
User hasn’t liked
367
Access Pattern Index Parameters Notes
WatchCategory
For User Main table N/A
Transaction to
increment Category
WatchCount and
ensure User hasn’t
watched
View Messages
for User Main table • Username Query to find all
Messages
View Unread
Messages for User GSI1
• Username Query to find all
Messages
Mark Message as
Read Main table N/A
Update Status and
remove GSI1
attributes
Send Hot New
Deal Message to
all Users
User index N/A 1. Scan UserIndex to
find all Users
Main table N/A 2. Create Message for
each User in step 1
Send new Brand
Deal Message to
all Brand
Watchers
Main table • BrandName
1.Query
BrandWatchers
partition to find
watchers for Brand
Main table N/A 2. Create Message for
each User in step 1
Send new
Category Deal
Message to all
Category
Watchers
Main table • CategoryName
1.Query
CategoryWatchers
partition to find
watchers for Category
Main table N/A 2. Create Message for
each User in step 1


