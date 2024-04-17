+++
title = 'DynamoDB'
date = '2024-04-12'
tags = ['aws']
subtitle = 'Some notes from the The DynamoDB Book'
+++

## Concepts

### Tables, Items & Attributes
- A table is a collection of items, each composed of one or more attributes.
- An attribute is given a type when it is written:
  - *Scalar*: single simple value (string, number, binary, boolean, null).
  - *Complex*: groupings with arbitrary nested attributes (list, map).
  - *Set*: multiple unique values of the same type (string sets, number sets, binary sets).
- The attribute type affects which operations can be performed on that attribute.
- Attributes with the same underlying value by a different type are not equal.
- A table is schemaless and attributes are not required on every item.
- Multiple different types of entities are commonly stored in a single table to allow for handling complex access patterns in a single request.
- Use as few tables as possible, ideally just one.

![items.png](/img/dynamodb/items.png)

### Primary Keys
- A primary key must be declared when creating a table.
- Each item must include the primary key and is uniquely identified by that primary key.
- There are two types of primary keys:
  - *Simple*: consists of a single partition key.
  - *Composite*: consists of a partition key and sort key.
- The terms "hash key" and "range key" are also used to refer to the partition key and sort key, respectively.
- A simple primary key allows for only fetching a single item at a time while a composite primary key enables fetching all items with the same partition key with *Query*.

### Secondary Indexes
- Secondary indexes allow for reshaping data into another format for querying.
- A key schema consisting of a partition key and (optional) sort key must be declared when creating a secondary index.
- All items will be copied from the base table into the secondary index in the reshaped form.
- There are two types of secondary indexes:
  - *Local*: Uses the same partition key as the primary key but a different sort key.
  - *Global*: Uses any attributes for the partition key and sort key.
- Filtering is built into the data model as the primary keys and secondary indexes determine how data is retrieved.

| | LSI | GSI |
|-|-|-|
| Throughput | Shared with base table | Separately provisioned |
| Consistency | Allow opting into strongly-consistent reads | Eventually-consistent reads only |
| Creation time | Must be specified when table created | Can be created and deleted as needed |

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
- When an item is written to the base table, it will be copied into the secondary index if it has the elements of the key schema for the secondary index.
- If the item does not have those elements then it will not be copied into the secondary index.
- A sparse index is a secondary index that intentionally excludes certain items from the base table to help satisfy certain access patterns.

### Item Collections
- An item collection consists of the group of items that share the same partition key in either the base table or a secondary index.
- All the items in an item collection will be allocated to the same partition.
- An item collection is ordered and stored as a B-tree.
- All items in an LSI are included as part of the same item collection as the base table.
- The item collection size includes both the size of the items in the base table and the size of the items in the local secondary index.

![item-collections.png](/img/dynamodb/item-collections.png)

### Partitioning
- Data is sharded across multiple partitions based on the partition key.
- This allows for horizontal scaling by adding more storage nodes.
- There are three storage nodes for each partition:
  - The primary node holds the authorative data.
  - Two secondary node provide durability and serve read requests.
- The request router serves as the frontend for all requests.
- An incoming write request is handled as follows:
  1. The request router uses the hash of the partition key to route the request to the appropriate primary node.
  2. The primary node commits the write and also commits the write to one of the two secondary nodes.
  3. The primary node responds to the client to indicate that the write was successful.
  4. The primary node asynchronously replicates the write to the third storage node.

### Consistency
- Consistency refers to whether a particular read operation receives all prior write operations.
- There are two consistency options available:
  - Strong: Any item read will reflect all prior writes.
  - Eventual: It is possible that read items will not reflect all prior writes.
- Reads are eventually-consistent by default thought it is possible to opt into strongly-consistent for base tables and LSIs.
- An eventually-consistent read consumes half the write capacity of a strongly-consistent read.
- Sources of eventual consistency include asynchronous data replication from primary to secondary nodes and from base tables to GSIs.

### Key Overloading 
- Different types of entities are often included in a single table.
- Key overloading refers to using generic names for the primary keys and using different values dependending on the item type.
- Prefixes are used for the partition and sort key values in order to identify the item type and to avoid overlap between different item types.
- Secondary indexes can be overloaded just like primary keys.
- Common generic names are: PK, SK, GSI1PK, GSI1SK, etc.

![key-overloading.png](/img/dynamodb/key-overloading.png)

### Time-to-Live
- TTLs allow for automatically deleting items after a specified time.
- To use TTLs, an attribute is specified on the table that will serve as the marker for item deletion.
- The attribute must be of type number.
- To expire an item, a Unix timestamp at seconds granularity can be stored in the specified attribute that indicates that time after which the item should be deleted.
- Items are usually deleted within 48 hours after the time indicated by the attribute.

### Capacity Units
- A RCU represents one strongly-consistent read per second or two eventually-consistent reads per second for an item up to 4KB.
- Transactional read requests require two RCUs to perform one read per second for items up to 4KB.
- Reading an item larger than 4KB will consume additional RCUs.
- One WCU represnts one write per second for an item up to 1KB in size.
- Transactional write requests require two WCUs to perform one write per second for items up to 1KB.
- Writing an item larger than 1KB will consume additional WCUs.

### Limits
- A single item is limited to 400KB.
- The *Query* and *Scan* actions will read a maximum of 1MB befor paginating (applied before any filter expressions).
- A single partition can have a maximum of 3000 RCUs or 1000 WCUs.
- When using a LSI, a single item collection cannot be larger than 10GB.
- There is no limit to the number of items that can be stored in a table.

### Efficiency
- By specifying the partition key, an operation starts with an *O(1)* lookup that reduces the dataset down to a maximum of 10GB on a particular storage node.
- For a *Query* on an item collection of size *n*, an *O(log n)* operation is required to find the starting value for the sort key after which a sequential read is limited to 1MB.

### Sorting
- Items need to be arranged so that they are sorted in advance.
- If a specific ordering is required when retrieving multiple items then a composite primary key must be used such that the ordering is performed with the sort key.
- Sort keys of type string or binary are sorted in order of UTF-8 bytes.
- Given uppercase letters are sorted before lowercase letters, sort keys should be standardized to a single case to avoid unexpected behavior.
- For timestamps to be sortable they should use a sortable format like a Unix epoch or ISO-8601.

## Operations

### *-Item
- Item-based actions operate on specific items: *GetItem*, *PutItem*, *UpdateItem* and *DeleteItem*.
- These are the only actions that can be used to write, update or delete items.
- The full primary key must be specified and the operation must be performed on the base table.
- *PutItem* can overwrite an existing item with the same primary key.
- *UpdateItem* will create an item if it does not exist.
- *UpdateItem* will only alter the properties specified and any other existing attributes will remain the same.

### Batch/Transaction
- Batch and transaction actions are used for operating on multiple items in a single request.
- These actions must specify the full primary key of items.
- Batch actions allow for reads or writes to succeed or fail independently.
- With transaction actions, the failure of a single operation will cause all the other writes to be rolled back.

### Query
- The *Query* action can be used to retrieve a contiguous block of items within a single item collection.
- The partition key must be specified and various conditions can be specified on the sort key: *>=*, *<=*, *begins_with()*, or *BETWEEN*.
- This can be performed against the base table or a secondary index.

![query-action.png](/img/dynamodb/query-action.png)

### Scan
- A *Scan* will retrieve all items in a table.
- In exceptional situations, a sparse secondary index can be modeled in a way that expects a *Scan*.
- Two optional properties are available to enable parallel scans:
  - *TotalSegments*: The total number of segments to split the *Scan* across.
  - *Segment*: The segment number to be scanned by this particular request.

### Optional Properties
- *ConsistentRead* is used to opt into strongly-consistent reads using *GetItem*, *BatchGetItem*, *Query* and *Scan*.
- *ScanIndexForward* controls which way a *Query* will read results from the sort key.
- *ReturnValues* determines which values a write operation returns:
  - *ALL_OLD*: Return all the attributes from the item before the operation was applied.
  - *UPDATED_OLD*: For any attributes updated in the operation, return the attributes before the operation was applied.
  - *ALL_NEW*: Return all the attributes from the item after the operation was applied.
  - *UPDATED_NEW*: For any attributes updated in the operation, return the attributes after the operation was applied.
- *ReturnConsumedCapacity* returns data about the capacity units that were consumed by the request.
- *ReturnItemCollectionMetrics* returns item collection size metrics.

![scan-index-forward.png](/img/dynamodb/scan-index-forward.png)

## Expressions

### Key Condition
- A key condition expression is used in a *Query* to describe which items to retrieve.
- It can only reference elements of the primary key.
- Sort key conditions can use simple comparisons: *>*, *<*, *=*, *begins_with*, *BETWEEN*.
- Every condition on the sort key can be expressed with the *BETWEEN* operator.
```py
client.query(
  TableName='CustomerOrders',
  KeyConditionExpression='#c = :c AND #ot BETWEEN :start and :end',
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
```py
client.query(
  TableName='MovieRoles',
  KeyConditionExpression='#actor = :actor',
  FilterExpression='#genre = :genre'
  ExpressionAttributeNames={
    '#actor': 'Actor',
    '#genre': 'Genre'
  },
  ExpressionAttributeValues={
    ':actor': { 'S': 'Tom Hanks' },
    ':genre': { 'S': 'Drama' }
  }
)
```

### Projection
- A project expression can be used in read operations to describe which attributes to return on read items.
- This can be used to access nested properties in a list or map attribute.
- Projection expressions are evaluated after the items are read from the table and the 1MB limit is reached.
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
- A condition expression is used in write operation to assert the existing condition of an item before writing to it.
- The operation will be canceled if the assertion fails.
- Condition expressions can operate on any attribute on the item.
- Several operators can be used: *>*, *<*, *=*, *BETWEEN*, *attribute_exists()*, *attribute_not_exists()*, *attribute_type()*, *begins_with()*, *contains()* and *size()*.
- The *TransactWriteItem* action can specify combination of different write operations (PutItem, UpdateItem or DeleteItem) or ConditionChecks.
```py
dynamodb.transact_write_items(
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
- Any combination of these verbs may be used in a single update statement.
- Multiple operations can be performed for a single verb.
```py
dynamodb.update_item(
  TableName='Users',
  Key={
    'Username': { 'S': 'python_fan' }
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

#### 1. Denormalization using a complex attribute
- Use an attribute with a complex data type like a list or map, violating NF1.
- It will not be possible to support access patterns based on the values in the complex attribute because a complex attribute cannot be used in a primary key.
- The amount of data in the complex attribute cannot be unbounded.

#### 2. Denormalization by duplicating data
- Duplicate the parent fields in each of the child items.
- This is particularly appropriate if the duplicated data is immutable.
- It may also be acceptable to duplicate the data depending on how often it changes and how many items contain the duplicated data.
- This balances the benefit of duplication (in the form of faster reads) against the costs of updating the data.

![one-to-many-1.png](/img/dynamodb/one-to-many-1.png)

#### 3. Composite primary key with Query
- Use a composite primary key and a *Query* to fetch multiple items within a single item collection.
- This solves for four common access patterns:
  1. Retrieve the parent item using *GetItem* and the primary key of the parent item.
  2. Retrieve the parent item and all its child items using *Query* and the partition key.
  3. Retrieve only the child items using *Query* with `begins_with(SK, "CHILD#")`.
  4. Retrieve a specific child item using *GetItem* and the primary key of the child item.

![one-to-many-2.png](/img/dynamodb/one-to-many-2.png)

#### 4. Secondary index with Query
- A *Query* can also be used with a secondary index if the primary key is reserved for some other purpose.
- This could be due to storing hierarchical data with a number of levels.

![one-to-many-3.png](/img/dynamodb/one-to-many-3.png)

![one-to-many-4.png](/img/dynamodb/one-to-many-4.png)

#### 5. Composite sort keys with hierarchical data
- It is not feasible to keep adding secondary indexes to enable arbitrary levels of fetching throughout a data hierarchy.
- A composite sort key refers to combining multiple properties together in the sort key to enable different search granularity.
- With each level separated by a hash in the sort key, we can search at different levels of granularity using `begins_with()`.
- This works well when there are more than two levels of hierarchy and access patterns for different levels within the hierarchy and when we want to return all sub-items in a level of the hierarchy rather than just the items in that level.

![one-to-many-5.png](/img/dynamodb/one-to-many-5.png)

### Many-to-Many

TODO

### Filtering

#### 1. Assemble Different Item Collections
- The key difference between this pattern and the simple filtering pattern with the sort key is that there’s no inherent meaning in the sort key values. Rather, the way that I’m sorting is a function of how I decided to arrange my items within a particular item collection.
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

![filtering-1.png](/img/dynamodb/filtering-1.png)

#### 2. Composite Sort Key
- A composite sort key is when you combine multiple data values in a sort key that allow you to filter on both values.
- This pattern works well when you always want to filter on two or more attributes in particular access pattern and one of the attributes is an enum-like value.
- They are sorted first by the OrderStatus, then by the OrderDate. This means we can do an exact match on that value and use more fine-grained filtering on the second value. This pattern would not work in reverse.
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
![filtering-2.png](/img/dynamodb/filtering-2.png)

#### 3. Sparse Indexes








This shows up most clearly in two situations:
• Filtering within an entity type based on a particular condition
• Projecting a single type of entity into a secondary index.


13.4.1. Using sparse indexes to provide a global filter on an item type

The first example of using a sparse index is when you filter within
an entity type based on a particular condition.
Imagine you had a SaaS application. In your table, you have
Organization items that represent the purchasers of your product.
Each Organization is made up of Members which are people that
have access to a given Organization’s subscription.
Your table might look as follows:
223
Notice that our Users have different roles in our application. Some
are Admins and some are regular Members. Imagine that you had
an access pattern that wanted to fetch all Users that had
Administrator privileges within a particular Organization.
If an Organization had a large number of Users and the condition
we want is sufficiently rare, it would be very wasteful to read all
Users and filter out those that are not administrators. The access
pattern would be slow, and we would expend a lot of read capacity
on discarded items.
Instead, we could use a sparse index to help. To do this, we would
add an attribute to only those User items which have Administrator
privileges in their Organization.
To handle this, we’ll add GSI1PK and GSI1SK attributes to
Organizations and Users. For Organizations, we’ll use
ORG#<OrgName> for both attributes. For Users, we’ll use
224
ORG#<OrgName> as the GSI1PK, and we’ll include Admin as the GSI1SK
only if the user is an administrator in the organization.
Our table would look as follows:
Notice that both Warren Buffett and Sheryl Sandberg have values
for GSI1SK but Charlie Munger does not, as he is not an admin.
Let’s take a look at our secondary index:
Notice we have an overloaded secondary index and are handling
multiple access patterns. The Organization items are put into a
single partition for fast lookup of all Organizations, and the User
225
items are put into partitions to find Admins in an Organization.
The key to note here is that we’re intentionally using a sparse index
strategy to filter out User items that are not Administrators. We can
still use non-sparse index patterns for other entity types.
The next strategy is a bit different. Rather than using an overloaded
index, it uses a dedicated sparse index to handle a single type of
entity.
13.4.2. Using sparse indexes to project a single
type of entity
A second example of where I like to use a sparse index is if I want to
project a single type of entity into an index. Let’s see an example of
where this can be useful.
Imagine I have an e-commerce application. I have several different
entity types in my application, including Customers that make
purchases, Orders that indicate a particular purchase, and
InventoryItems that represent products I have available for sale.
My table might look as follows:
226
Notice that the table includes Customers, Orders, and
InventoryItems, as discussed, and these items are interspersed
across the table.
My marketing department occasionally wants to send marketing
emails to all Customers to alert them of hot sales or new products.
To find all my Customers in my base table is an expensive task, as I
would need to scan my entire table and filter out the items that
aren’t Customers. This is a big waste of time and of my table’s read
capacity.
Instead of doing that, I’ll add an attribute called CustomerIndexId
on my Customer items. Now my table looks as follows:
227
Notice the Customer items now have an attribute named
CustomerIndexId as outlined in red.
Then, I create a secondary index called CustomerIndex that uses
CustomerIndexId as the partition key. Only Customer items have
that attribute, so they are the only ones projected into that index.
The secondary index looks as follows:
Only the Customer items are projected into this table. Now when
the marketing department wants to find all Customers to send
marketing emails, they can run a Scan operation on the
CustomerIndex, which is much more targeted and efficient. By
isolating all items of a particular type in the index, our sparse index
makes finding all items of that type much faster.
228
Again, notice that this strategy does not work with index
overloading. With index overloading, we’re using a secondary index
to index different entity types in different ways. However, this
strategy relies on projecting only a single entity type into the
secondary index.

Both sparse index patterns are great for filtering out non-matching items entirely.

#### 6. Filter Expressions

- You can include filter expressions in Query and Scan operations to remove items from your results that don’t match a given condition.
- This is because a filter expression is applied after items are read, meaning you pay for all of the items that get filtered out and you are subject to the 1MB results limit before your filter is evaluated.
- Because of this, you cannot count on filter expressions to save a bad model. Filter expressions are, at best, a way to slightly improve the performance of a data model that already works well.
- Reducing response payload size.
- Easier application filtering.
- Better validation around time-to-live (TTL) expiry.
- If you implement this filter via a filter expression, you don’t know how many items you will need to fetch to ensure you get ten orders to return to the client. Accordingly, you’ll likely need to vastly overfetch your items or have cases where you make follow-up requests to retrieve additional items.

- This can be useful for reducing response payload size and to filter out items with expired TTLs.

#### 7. Client-side Filtering
- Delay filtering to the client.
- This is appropriate when filtering is difficult to model in the database or when the dataset is small.

### Sorting

#### 1. Sorting on changing attributes
- Including an *UpdatedAt* field in the sort key is undesirable as it would require deleting and re-creating the item whenever it is updated.
- Instead, use immutable values for the primary key and introduce a secondary index where the *UpdatedAt* is the sort key.

![sorting-1.png](/img/dynamodb/sorting-1.png)

#### 2. Ascending vs. Descending


As we discussed in Chapter 5, you can use the ScanIndexForward
property to tell DynamoDB how to order your items.

One complication arises when you are combining the one-to-many
relationship strategies from Chapter 9 with the sorting strategies in
this chapter.

With those strategies, you are often combining
multiple types of entities in a single item collection and using it to
read both a parent and multiple related entities in a single request.

When you do this, you need to consider the common sort order you’ll use in this access pattern to know where to place the parent item.

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


#### Faking ascending order

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

### Migrations
- A purely additive change means you can start writing the new attributes or items without changes to existing items.
- A purely additive change is much easier, while editing existing items usually requires a batch job to scan existing items and decorate them with new attributes.
- When considering whether a change is additive, you only need to consider indexing attributes, not application attributes.

#### Adding new attributes to an existing entity
- The first, and easiest, type of migration is to add new attributes to existing entities. When adding new attributes, you can simply add this in your application code.
- Adding application attributes is the easiest type of migration because you don’t really think about application attributes while modeling your data in DynamoDB. You are primarily concerned with the attributes that are used for indexing within DynamoDB.
- The schemaless nature of DynamoDB makes it simple to add new attributes in your application without doing large-scale migrations.

#### Adding a new entity type without relations
- If there are no access patterns that require modeling a relation with existing entity types then we can simply start writing the new entity type without making changes to existing items.
- If the new entity type needs to be added into an existing item collection 


This one is similar to the last in that you’re adding a new entity type
into your table, but it’s different in that you do have a relational
access pattern of "Fetch parent entity and its related entities".
In this scenario, you can try to reuse an existing item collection. This
is great if your parent entity is in an item collection that’s not being
used for an existing relationship.
259
As an example, think about an application like Facebook before
they introduced the "Like" button. Their DynamoDB table may
have had Post items that represented a particular post that a user
made.
The basic table might look as follows:
Then someone has the brilliant idea to allow users to "like" posts via
the Like button. When we want to add this feature, we have an
access pattern of "Fetch post and likes for post".
In this situation, the Like entity is a completely new entity type, and
we want to fetch it at the same time as fetching the Post entity. If we
look at our base table, the item collection for the Post entity isn’t
being used for anything. We can add our Like items into that
collection by using the following primary key pattern for Likes:
• PK: POST#<PostId>
• SK: LIKE#<Username>
When we add a few items to our table, it looks like this:
260
Outlined in red is an item collection with both a Post and the Likes
for that Post. Like our previous examples, this is a purely additive
change that didn’t require making changes to our existing Post
items. All that was required is that we modeled our data
intentionally to locate the Like items into the existing item
collection for Posts.


15.4. Adding a new entity type into a new item collection
We have a new item type and a relational access pattern with an existing type. However, there’s a twist—we don’t have an existing item collection where we can handle this relational access pattern.

Let’s continue to use the example from the last section. We have a social application with Posts and Likes. Now we want to add Comments. Users can comment on a Post to give encouragement or argue about some pedantic political point. With our new Comment entity, we have a relational access pattern where we want to fetch a Post and the most recent Comments for that Post. How can we model this new pattern? The Post item collection is already being used on the base table. To handle this access pattern, we’ll need to create a new item collection in a global secondary index. To do this, let’s add the following attributes to the Post item:

• GSI1PK: POST#<PostId>
• GSI1SK: POST#<PostId>

And we’ll create a Comment item with the following attributes:

• PK: COMMENT#<CommentId>
• SK: COMMENT#<CommentId>
• GSI1PK: POST#<PostId>
• GSI1SK: COMMENT#<Timestamp>

When we do that, our base table looks as follows:
Notice that we’ve added two Comment items at the bottom
outlined in blue. We’ve also add two attributes, GSI1PK and GSI1SK,
to our Post items, outlined in red .
When we switch to our secondary index, it looks as follows:
263
Our Post item is in the same item collection as our Comment items,
allowing us to handle our relational use case.
This looks easy in the data model, but I’ve skipped over the hard
part. How do you decorate the existing Post items to add GSI1PK
and GSI1SK attributes?
You will need to run a table scan on your table and update each of
the Post items to add these new attributes. A simplified version of
the code is as follows:
264
last_evaluated = ''
params = {
  "TableName": "SocialNetwork",
  "FilterExpression": "#type = :type",
  "ExpressionAttributeNames": {
  "#type": "Type"
  },
  "ExpressionAttributeValues": {
  ":type": { "S": "Post" }
  }
}
while True:
  if last_evaluated:
  params['ExclusiveStartKey'] = last_evaluated
  results = client.scan(**params)
  for item in results['Items']:
  client.update_item(
  TableName='SocialNetwork',
  Key={
  'PK': item['PK'],
  'SK': item['SK']
  },
  UpdateExpression="SET #gsi1pk = :gsi1pk, #gsi1sk = :gsi1sk",
  ExpressionAttributeNames={
  '#gsi1pk': 'GSI1PK',
  '#gsi1sk': 'GSI1SK'
  }
  ExpressionAttributeValues={
  ':gsi1pk': item['PK'],
  ':gsi1sk': item['SK']
  }
  )
  if not results['LastEvaluatedKey']:
  break
  last_evaluated = results['LastEvaluatedKey']
This script is running a Scan API action against our DynamoDB
table. It’s using a filter expression to filter out any items whose Type
attribute is not equal to Post, as we’re only updating the Post items
in this job.
As we receive items from our Scan result, we iterate over those
items and make an UpdateItem API request to add the relevant
265
properties to our existing items.
There’s some additional work to handle the LastEvaluatedKey
value that is received in a Scan response. This indicates whether we
have additional items to scan or if we’ve reached the end of the
table.
There are a few things you’d want to do to make this better,
including using parallel scans, adding error handling, and updating
multiple items in a BatchWriteItem request, but this is the general
shape of your ETL process. There is a note on parallel scans at the
end of this chapter.
This is the hardest part of a migration, and you’ll want to test your
code thoroughly and monitor the job carefully to ensure all goes
well. However, there’s really not that much going on. A lot of this
can be parameterized:
• How do I know which items I want?
• Once I get my items, what new attributes do I need to add?.
From there, you just need to take the time for the whole update
operation to run.

15.5. Joining existing items into a new
item collection
So far, all of the examples have involved adding a new item type
into our application. But what if we just have a new access pattern
on existing types? Perhaps we want to filter existing items in a
different way. Or maybe we want to use a different sorting
mechanism. We may even want to join two items that previously
were separate.

Find the items you
want to update and design an item collection in a new or existing
secondary index. Then, run your script to add the new attributes so
they’ll be added to the secondary index.

### Other

#### Ensuring uniqueness on two or more attributes
- To ensure that a particular attribute is unique it must be built directly into the primary key structure.
- A condition expression can be used to ensure uniqueness when writing.
- The combination of a partition key and sort key makes an item unique.
- To ensure that multiple attributes are unique across the table, it will be necessary to write both attributes in a transaction with each put asserting that the attribute does not exist 

In this example, we’ll create two items in a transaction: one that
tracks the user by username, and one that tracks the email by
username.
271
The code to write such a transaction would be as follows:
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
  ...
  },
  'ConditionExpression': 'attribute_not_exists(PK)
  }
  },
  {
  'Put': {
  'TableName': 'UsersTable',
  'Item': {
  'PK': { 'S': 'USEREMAIL#alex@debrie.com' },
  'SK': { 'S': 'USEREMAIL#alex@debrie.com' },
  },
  'ConditionExpression': 'attribute_not_exists(PK)
  }
  }
  ]
)
For each write operation, we’re including a condition expression
that ensures that an item with that primary key doesn’t exist. This
confirms that the username is not currently in use and that the
email address is not in use.
And now your table would look as follows:
272
Notice that the item that stores a user by email address doesn’t have
any of the user’s properties on it. You can do this if you will only
access a user by a username and never by an email address. The
email address item is essentially just a marker that tracks whether
the email has been used.
If you will access a user by email address, then you need to
duplicate all information across both items. Then your table might
look as follows:
I’d avoid this if possible. Now every update to the user item needs
to be a transaction to update both items. It will increase the cost of
your writes and the latency on your requests.


16.2. Handling sequential IDs
In relational database systems, you often use a sequential ID as a
primary key identifier for each row in a table. With DynamoDB,
this is not the case. You use meaningful identifiers, like usernames,
product names, etc., as unique identifiers for your items.
That said, sometimes there are user-facing reasons for using
sequential identifiers. Perhaps your users are creating entities, and
it’s easiest to assign them sequential identifiers to keep track of
them. Examples here are Jira tickets, GitHub issues, or order
numbers.
DynamoDB doesn’t have a built-in way to handle sequential
identifiers for new items, but you can make it work through a
combination of a few building blocks.
Let’s use a project tracking software like Jira as an example. In Jira,
you create Projects, and each Project has multiple Issues. Issues are
given a sequential number to identify them. The first issue in a
project would be Issue #1, the second issue would be Issue #2, etc.
When a new issue is created within a project, we’ll do a two-step
process.
First, we will run an UpdateItem operation on the Project item to
increment the IssueCount attribute by 1. We’ll also set the
ReturnValues parameter to UPDATED_NEW, which will return the
current value of all updated attributes in the operation. At the end
of this operation, we will know what number should be used for our
new issue. Then, we’ll create our new Issue item with the new issue
number.
274
The full operation looks as follows:
resp = client.update_item(
  TableName='JiraTable',
  Key={
  'PK': { 'S': 'PROJECT#my-project' },
  'SK': { 'S': 'PROJECT#my-project' }
  },
  UpdateExpression="SET #count = #count + :incr",
  ExpressionAttributeNames={
  "#count": "IssueCount",
  },
  ExpressionAttributeValues={
  ":incr": { "N": "1" }
  },
  ReturnValues='UPDATED_NEW'
)
current_count = resp['Attributes']['IssueCount']['N']
resp = client.put_item(
  TableName='JiraTable',
  Item={
  'PK': { 'S': 'PROJECT#my-project' },
  'SK': { 'S': f"ISSUE#{current_count}" },
  'IssueTitle': { 'S': 'Build DynamoDB data model' }
  ... other attributes ...
  }
)
First we increment the IssueCount on our Project item. Then we
use the updated value in the PutItem operation to create our Issue
item.
This isn’t the best since you’re making two requests to DynamoDB
in a single access pattern. However, it can be a way to handle autoincrementing IDs when you need them.
16.3. Pagination
Pagination is a common requirement for APIs. A request will often
return the first 10 or so results, and you can receive additional
results by making a follow-up request and specifying the page.
275
In a relational database, you may use a combination of OFFSET and
LIMIT to handle pagination. DynamoDB does pagination a little
differently, but it’s pretty straightforward.
When talking about pagination with DynamoDB, you’re usually
paginating within a single item collection. You’re likely doing a
Query operation within an item collection to fetch multiple items
at once.
Let’s walk through this with an example. Imagine you have an ecommerce application where users make multiple orders. Your
table might look like this:
Our table includes Order items. Each Order item uses a PK of
USER#<username> for grouping and then uses an OrderId in the sort
key. The OrderId is a KSUID (see the notes in Chapter 14 for notes
on KSUIDs) which gives them rough chronological ordering.
A key access pattern is to fetch the most recent orders for a user.
We only return five items per request, which seems very small but
it makes this example must easier to display. On the first request,
we would write a Query as follows:
276
resp = client.query(
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
This request uses our partition key to find the proper partition, and
it looks for all values where the sort key is less than ORDER$. This
will come immediately after our Order items, which all start with
ORDER#. Then, we use ScanIndexForward=False to read items
backward (reverse chronological order) and set a limit of 5 items.
It would retrieve the items outlined in red below:
This works for the first page of items but what if a user makes
follow-up requests?
In making that follow-up request, the URL would be something like
277
https://my-ecommercestore.com/users/alexdebrie/orders?before=1YRfXS14inXwIJEf9
tO5hWnL2pi. Notice how it includes both the username (in the URL
path) as well as the last seen OrderId (in a query parameter).
We can use these values to update our Query:
resp = client.query(
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
Notice how the SK value we’re comparing to is now the Order item
that we last saw. This will get us the next item in the item collection,
returning the following result:
By building these hints into your URL, you can discover where to
278
start for your next page through your item collection.
16.4. Singleton items
In most of our examples, we create an item pattern that is reused
for the primary key across multiple kinds of items.
If you’re creating a User item, the pattern might be:
• PK: USER#<Username>
• SK: USER#<Username>
If you’re creating an Order item for the user, the pattern might be:
• PK: USER#<Username>
• SK: ORDER#<OrderId>
Notice that each item is customized for the particular user or order
by virtue of the username and/or order id.
But sometimes you don’t need a customized item. Sometimes you
need a single item that applies across your entire application.
For example, maybe you have a table that is tracking some
background jobs that are running in your application. Across the
entire application, you want to have at most 100 jobs running at any
given time.
To track this, you could create a singleton item that is responsible for
tracking all jobs in progress across the application. Your table might
look as follows:
279
We have three Job items outlined in red. Further, there is a
singleton item with a PK and SK of JOBS that tracks the existing jobs
in progress via a JobsInProgress attribute of type string set. When
we want to start a new job, we would do a transaction with two write
operations:
1. Update the JobsInProgress attribute of the singleton Jobs item
to add the new Job ID to the set if the current length of the set is
less than 100, the max number of jobs in progress.
2. Update the relevant Job item to set the status to IN_PROGRESS.
Another example for using singleton items is in the Big Time Deals
example in Chapter 20 where we use singleton items for the front
page of the application.
16.5. Reference counts
We talked about one-to-many relationships and many-to-many
relationships in previous chapters. These are very common
280
patterns in your application.
Often you’ll have a parent item in your application that has a large
number of related items through a relationship. It could be the
number of retweets that a tweet receives or the number of stars for
a GitHub repo.
As you show that parent item in the UI, you may want to show a
reference count of the number of related items for the parent. We
want to show the total number of retweets and likes for the tweet
and the total number of stars for the GitHub repo.
We could query and count the related items each time, but that
would be highly inefficient. A parent could have thousands or more
related items, and this would burn a lot of read capacity units just to
receive an integer indicating the total.
Instead, we’ll keep a count of these related items as we go.
Whenever we add a related item, we usually want to do two things:
1. Ensure the related item doesn’t already exist (e.g. this particular
user hasn’t already starred this repo);
2. Increase the reference count on the parent item.
Note that we only want to allow it to proceed if both portions
succeed. This is a great case for DynamoDB Transactions!
Remember that DynamoDB Transactions allow you to combine
multiple operations in a single request, and the operations will only
be applied if all operations succeed.
The code to handle this would look as follows:
281
result = dynamodb.transact_write_items(
  TransactItems=[
  {
  "Put": {
  "Item": {
  "PK": { "S": "REPO#alexdebrie#dynamodb-book" },
  "SK": { "S": "STAR#danny-developer" }
  ...rest of attributes ...
  },
  "TableName": "GitHubModel",
  "ConditionExpression": "attribute_not_exists(PK)"
  }
  },
  {
  "Update": {
  "Key": {
  "PK": { "S": "REPO#alexdebrie#dynamodb-book" },
  "SK": { "S": "#REPO#alexdebrie#dynamodb-book" }
  },
  "TableName": "GitHubModel",
  "ConditionExpression": "attribute_exists(PK)",
  "UpdateExpression": "SET #count = #count + :incr",
  "ExpressionAttributeNames": {
  "#count": "StarCount"
  },
  "ExpressionAttributeValues": {
  ":incr": { "N": "1" }
  }
  }
  }
  ]
)
Notice that in our transaction, we’re doing a PutItem operation to
insert the Star entity to track this user starring a repo. We include a
condition expression to ensure the user hasn’t already starred the
repo. Additionally, we have an UpdateItem operation that
increments the StarCount attribute on our parent Repo item.
I use this reference count strategy quite a bit in my applications,
and DynamoDB Transactions have made it much easier than the
previous mechanism of editing multiple items and applying
manual rollbacks in the event of failure.
282
16.6. Conclusion
In this chapter we covered additional strategies for DynamoDB.
The strategies are summarized below.
Strategy Notes Relevant examples
Ensuring uniqueness on
two or more attributes
Create a tracking item
and use DynamoDB
Transactions
E-commerce example
(Users with unique email
and username)
Handling sequential IDs Track current ID in
attribute. Increment that
attribute and use
returned value to create
new item.
GitHub Issues & Pull
Requests
Pagination Build pagination
elements into URL
structure
E-commerce orders
Singleton items Use for tracking global
state or for assembling a
meta view
Deals example
Reference counts Use a transaction to
maintain a count of
related items on a parent
item
Deals example; GitHub
example
Table 13. Additional strategies

## Examples

1. Create an entity-relationship diagram
2. Define your access patterns
3. Model your primary key structure
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






Entity PK SK
Customers CUSTOMER#<Username> CUSTOMER#<Username>
CustomerEmails CUSTOMEREMAIL#<Email> CUSTOMEREMAIL#<Email>
Addresses
Orders
OrderItems


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


