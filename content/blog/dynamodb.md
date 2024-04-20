+++
title = "DynamoDB"
date = "2024-04-12"
tags = ["aws"]
subtitle = "Some notes from the The DynamoDB Book"
+++

<!--
Locked in: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
Done: 13
In-progress: 14
Omitted: 12, 15
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

#### 1. Sorting on changing attributes
- Including an *UpdatedAt* field in the sort key is undesirable as it would require deleting and re-creating the item whenever it is updated.
- Instead, use immutable values for the primary key and introduce a secondary index where the *UpdatedAt* is the sort key.

![sorting-1.png](/img/dynamodb/sorting-1.png)

<!--
Like with joins and filtering, you need to arrange your items so
they’re sorted in advance.
If you need to have specific ordering when retrieving multiple
items in DynamoDB, there are two main rules you need to follow.
First, you must use a composite primary key. Second, all ordering
234
must be done with the sort key of a particular item collection.
Think back to our discussion in Chapter 4 about how DynamoDB
enforces efficiency. First, it uses the partition key to isolate item
collections into different partitions and enables an O(1) lookup to
find the proper node. Then, items within an item collection are
stored as a B-tree which allow for O(log n) time complexity on
search. This B-tree is arranged in lexicographical order according
to the sort key, and it’s what you’ll be using for sorting.
In this chapter, we’ll review some strategies for sorting. We’ll cover:
• Basics of sorting (what is lexicographical sorting, how to handle
timestamps, etc.)
• Sorting on changing attributes
• Ascending vs. descending
• Two one-to-many access patterns in a single item collection
• Zero-padding with numbers
• Faking ascending order
It’s a lot to cover, so let’s get started!
14.1. Basics of sorting
Before we dive into actual sorting strategies, I want to cover the
basics of how DynamoDB sorts items.
As mentioned, sorting happens only on the sort key. You can only
use the scalar types of string, number, and binary for a sort key.
Thus, we don’t need to think about how DynamoDB would sort a
map attribute!
235
For sort keys of type number, the sorting is exactly as you would
expect—items are sorted according to the value of the number.
For sort keys of type string or binary, they’re sorted in order of
UTF-8 bytes. Let’s take a deeper look at what that means.
14.1.1. Lexicographical sorting
A simplified version of sorting on UTF-8 bytes is to say the
ordering is lexicographical. This order is basically dictionary order
with two caveats:
1. All uppercase letters come before lowercase letters
2. Numbers and symbols (e.g. # or $) are relevant too.
The biggest place I see people get tripped up with lexicographical
ordering is by forgetting about the uppercase rule. For example, my
last name is "DeBrie" (note the capital "B" in the middle). Imagine
you had Jimmy Dean, Laura Dern, and me in an item collection
using our last names. If you forgot about capitalization, it might
turn out as follows:
You might be surprised to see that DeBrie came before Dean! This
is due to the casing—uppercase before lowercase.
236
To avoid odd behavior around this, you should standardize your
sort keys in all uppercase or all lowercase values:
With all last names in uppercase, they are now sorted as we would
expect. You can then hold the properly-capitalized value in a
different attribute in your item.
14.1.2. Sorting with Timestamps
A second basic point I want to cover with sorting is how to handle
timestamps. I often get asked the best format to use for timestamps.
Should we use an epoch timestamp (e.g. 1583507655) with
DynamoDB’s number type? Or should we use IS0-8601 (e.g. 2020-
03-06T15:14:15)? Or something else?
First off, your choice needs to be sortable. In this case, either epoch
timestamps or ISO-8601 will do. What you absolutely cannot do is
use something that’s not sortable, such as a display-friendly format
like "May 26, 1988". This won’t be sortable in DynamoDB, and you’ll
be in a world of hurt.
Beyond that, it doesn’t make a huge difference. I prefer to use ISO8601 timestamps because they’re human-readable if you’re
debugging items in the DynamoDB console. That said, it can be
237
tough to decipher items in the DynamoDB console if you have a
single-table design. As mentioned in Chapter 9, you should have
some scripts to aid in pulling items that you need for debugging.
14.1.3. Unique, sortable IDs
A common need is to have unique, sortable IDs. This comes up
when you need a unique identifier for an item (and ideally a
mechanism that’s URL-friendly) but you also want to be able to sort
a group of these items chronologically. This problem comes up in
both the Deals example (Chapter 17) and the GitHub Migration
example (Chapter 19).
There are a few options in this space, but I prefer the KSUID
implementation from the folks at Segment. A KSUID is a KSortable Unique Identifier. Basically, it’s a unique identifier that is
prefixed with a timestamp but also contains enough randomness to
make collisions very unlikely. In total, you get a 27-character string
that is more unique than a UUIDv4 while still retaining
lexicographical sorting.
Segment released a CLI tool for creating and inspecting KSUIDs.
The output is as follows:
ksuid -f inspect
REPRESENTATION:
  String: 1YnlHOfSSk3DhX4BR6lMAceAo1V
  Raw: 0AF14D665D6068ACBE766CF717E210D69C94D115
COMPONENTS:
  Time: 2020-03-07T13:02:30.000Z
  Timestamp: 183586150
  Payload: 5D6068ACBE766CF717E210D69C94D115
The "String" version of it (1YnlHOfSSk3DhX4BR6lMAceAo1V) shows
the actual value you would use in your application. The various
238
components below show the time and random payload that were
used.
There are implementations of KSUIDs for many of the popular
programming languages, and the algorithm is pretty
straightforward if you do need to implement it yourself.

Thanks to Rick Branson and the folks at Segment for the implementation and
description of KSUIDs. For more on this, check out Rick’s blog post on the
implementation of KSUIDs.
14.2. Sorting on changing attributes
The sort key in DynamoDB is used for sorting items within a given
item collection. This can be great for a number of purposes,
including viewing the most recently updated items or a leaderboard
of top scores. However, it can be tricky if the value you are sorting
on is frequently changing. Let’s see this with an example.
Imagine you have a ticket tracking application. Organizations sign
up for your application and create tickets. One of the access
patterns is to allow users to view tickets in order by the most
recently updated.
You decide to model your table as follows:
239
In this table design, the organization name is the partition key,
which gives us the ‘group by’ functionality. Then the timestamp for
when the ticket was last updated is the sort key, which gives us
‘order by’ functionality. With this design, we could use
DynamoDB’s Query API to fetch the most recent tickets for an
organization.
However, this design causes some problems. When updating an
item in DynamoDB, you may not change any elements of the
primary key. In this case, your primary key includes the UpdatedAt
field, which changes whenever you update a ticket. Thus, anytime
you update a ticket item, we would need first to delete the existing
ticket item, then create a new ticket item with the updated primary
key.
We have caused a needlessly complicated operation and one that
could result in data loss if you don’t handle your operations
correctly.
240
Instead, let’s try a different approach. For our primary key, let’s use
two attributes that won’t change. We’ll keep the organization name
as the partition key but switch to using TicketId as the sort key.
Now our table looks as follows:
Now we can add a secondary index where the partition key is
OrgName and the sort key is UpdatedAt. Each item from the base
table is copied into the secondary index, and it looks as follows:
241
Notice that this is precisely how our original table design used to
look. We can use the Query API against our secondary index to
satisfy our ‘Fetch most recently updated tickets’ access pattern.
More importantly, we don’t need to worry about complicated
delete + create logic when updating an item. We can rely on
DynamoDB to handle that logic when replicating the data into a
secondary index.
14.3. Ascending vs. descending
Now that we’ve got the basics out of the way, let’s look at some
more advanced strategies.
As we discussed in Chapter 5, you can use the ScanIndexForward
property to tell DynamoDB how to order your items. By default,
DynamoDB will read items in ascending order. If you’re working
with words, this means starting at aardvark and going toward zebra.
242
If you’re working with timestamps, this means starting at the year
1900 and working toward the year 2020.
You can flip this by using ScanIndexForward=False, which means
you’ll be reading items in descending order. This is useful for a
number of occasions, such as when you want to get the most recent
timestamps or you want to find the highest scores on the
leaderboard.
One complication arises when you are combining the one-to-many
relationship strategies from Chapter 9 with the sorting strategies in
this chapter. With those strategies, you are often combining
multiple types of entities in a single item collection and using it to
read both a parent and multiple related entities in a single request.
When you do this, you need to consider the common sort order
you’ll use in this access pattern to know where to place the parent
item.
For example, imagine you have an IoT device that is sending back
occasional sensor readings. One of your common access patterns is
to fetch the Device item and the most recent 10 Reading items for
the device.
You could model your table as follows:
243
Notice that the parent Device item is located before any of the
Reading items because "DEVICE" comes before "READING" in the
alphabet. Because of this, our Query to get the Device and the
Readings would retrieve the oldest items. If our item collection was
big, we might need to make multiple pagination requests to get the
most recent items.
To avoid this, we can add a # prefix to our Reading items. When we
do that, our table looks as follows:
244
Now we can use the Query API to fetch the Device item and the
most recent Reading items by starting at the end of our item
collection and using the ScanIndexForward=False property.
When you are co-locating items for one-to-many or many-tomany relationships, be sure to consider the order in which you
want the related items returned so that your parent itself is located
accordingly.
14.4. Two relational access patterns in a
single item collection
Let’s take that last example up another level. If we can handle a
one-to-many relationship where the related items are located
before or after the parent items, can we do both in one item
collection?
245
We sure can! To do this, you’ll need to have one access pattern in
which you fetch the related items in ascending order and another
access pattern where you fetch the related items in descending
order. It also works if order doesn’t really matter for one of the two
access patterns.
For example, imagine you had a SaaS application. Organizations
sign up and pay for your application. Within an Organization, there
are two sub-concepts: Users and Teams. Both Users and Teams
have a one-to-many relationship with Organizations.
For each of these relationships, imagine you had a relational access
pattern of "Fetch Organization and all {Users|Teams} for the
Organization". Further, you want to fetch Users in alphabetical
order, but Teams can be returned in any order because there
shouldn’t be that many of them.
You could model your table as follows:
In this table, we have our three types of items. All share the same PK
value of ORG#<OrgName>. The Team items have a SK of
246
#TEAM#<TeamName>. The Org items have an SK of ORG#<OrgName>.
The User items have an SK of USER#<UserName>.
Notice that the Org item is right between the Team items and the
User items. We had to specifically structure it this way using a #
prefix to put the Team item ahead of the Org item in our item
collection.
Now we could fetch the Org and all the Team items with the
following Query:
result = dynamodb.query(
  TableName='SaaSTable',
  KeyConditionExpression="#pk = :pk AND #sk <= :sk",
  ExpressionAttributeNames={
  "#pk": "PK",
  "#sk": "SK"
  },
  ExpressionAttributeValues={
  ":pk": { "S": "ORG#MCDONALDS" },
  ":sk": { "S": "ORG#MCDONALDS" }
  },
  ScanIndexForward=False
)
This goes to our partition and finds all items less than or equal to
than the sort key value for our Org item. Then it scans backward to
pick up all the Team items.
Below is an image of how it works on the table:
247
You can do the reverse to fetch the Org item and all User items:
look for all items greater than or equal to our Org item sort key,
then read forward to pick up all the User items.
This is a more advanced pattern that is by no means necessary, but
it will save you additional secondary indexes in your table. You can
see a pattern of this in action in Chapter 21.
14.5. Zero-padding with numbers
You may occasionally want to order your items numerically even
when the sort key type is a string. A common example of this is
when you’re using prefixes to indicate your item types. Your sort
key might be <ItemType>#<Number>. While this is doable, you need
to be careful about lexicographic sorting with numbers in strings.
As usual, this is best demonstrated with an example. Imagine in our
IoT example above that, instead of using timestamps in the sort
key, we used a reading number. Each device kept track of what
248
reading number it was on and sent that up with the sensor’s value.
You might have a table that looks as follows:
Yikes, our readings are out of order. Reading number 10 is placed
ahead of Reading number 2. This is because lexicographic sorting
evalutes one character at a time, from left to right. When it is
compared "10" to "2", the first digit of 10 ("1") is before the first digit
of 2 ("2"), so 10 was placed before 2.
To avoid this, you can zero-pad your numbers. To do this, you
make the number part of your sort key a fixed length, such as 5
digits. If your value doesn’t have that many digits, you use a "0" to
make it longer. Thus, "10" becomes "00010" and "2" becomes
"00002".
Now our table looks as follows:
249
Now Reading number 2 is placed before Reading number 10, as
expected.
The big factor here is to make sure your padding is big enough to
account for any growth. In our example, we used a fixed length of 5
digits which means we can go to Reading number 99999. If your
needs are bigger than that, make your fixed length longer.
I’d recommend going to the maximum number of related items
you could ever imagine someone having, then adding 2-3 digits
beyond that. The cost is just an extra few characters in your item,
but the cost of underestimating will add complexity to your data
model later on. You may also want to have an alert condition in
your application code that lets you know if a particular count gets
to more than X% of your maximum, where X is probably 30 or so.
14.6. Faking ascending order
This last sorting pattern is a combination of a few previous ones,
and it’s pretty wild. Imagine you had a parent entity that had two
250
one-to-many relationships that you want to query. Further,
imagine that both of those one-to-many relationships use a number
for identification but that you want to fetch both relationships in
the same order (either descending or ascending).
The problem here is that because you want to fetch the
relationships in the same order, you would need to use two
different item collections (and thus secondary indexes) to handle it.
One item collection would handle the parent entity and the most
recent related items of entity A, and the other item collection would
handle the parent entity and the most recent related items of entity
B.
However, we could actually get this in the same item collection
using a zero-padded difference. This is similar to our zero-padded
number, but we’re storing the difference between the highest
number available and the actual number of our item.
For example, imagine we were again using a zero-padded number
with a width of 5. If we had an item with an ID of "157", the zeropadded number would be "00157".
To find the zero-padded difference, we subtract our number ("157")
from the highest possible value ("99999"). Thus, the zero-padded
difference would be "99842".
If we put this in our table, our items would look as follows:
251
Notice that I changed the SK structure of the Reading items so that
the parent Device item is now at the top of our item collection. Now
we can fetch the Device and the most recent Readings by starting at
Device and reading forward, even though we’re actually getting the
readings in descending order according to their ReadingId.
This is a pretty wacky pattern, and it actually shows up in the
GitHub Migration example in Chapter 19. That said, you may not
ever have a need for this in practice. The best takeaway you can get
from this strategy is how flexible DynamoDB can be if you
combine multiple strategies. Once you learn the basics, you can
glue them together in unique ways to solve your problem.
-->

### Other

#### Ensuring uniqueness on two or more attributes
- To ensure that a particular attribute is unique it must be built directly into the primary key structure.
- A condition expression can be used to ensure uniqueness when writing.
- The combination of a partition key and sort key makes an item unique.
- To ensure that multiple attributes are unique across the table, it will be necessary to write both attributes in a transaction with each put asserting that the attribute does not exist 

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

