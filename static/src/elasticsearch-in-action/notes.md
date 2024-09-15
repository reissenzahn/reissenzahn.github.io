# Overview


# Concepts

## Ingestion
- Data can be indexed into Elasticsearch from multiple sources (databases, file systems, applications, etc.).
- An extract-transform-load tool is commonly used to transform and enrich data before it is indexed.
- Full-text fields undergo text analysis consisting of tokenization and normalization. The resulting tokens, root words and synonyms are stored in an inverted index.

## Documents
- A document is the basic unit of information indexed by Elasticsearch, represented as JSON.
- Mapping uses a schema definition to convert JSON data types into appropriate Elasticsearch data types.
- Data in Elasticsearch is de-normalized and a document consists of self-contained information with no relations.
- Documents can be indexed without a predefined schema.
- There are two types of document APIs: single-document and multiple-document.

## Indexes
- An index is a logical collection of documents composed of shards.
- For instance, a index might be composed of three shards on three nodes, one shard per node. In addition, it might have two replicas per shard, both hosted on other nodes.
- Each index should have only one document type. A single index may not contain fields of different data types.
- An index can hold any number of documents.
- By default, a index is created with a single shard and replica.

## Replicas
- Replicas provide data redundancy and serve reads.
- The number of replicas can be modified on a live index.

## Shards
- A shard is a running instance of Lucene that handles storage and retrieval of data.
- The number of shards cannot be modified on a live index.
- Index settings allow for configuring the shards and replicas.
- Each document must belong to a particular primary shard.

## Aliases
- Aliases are alternate names given to a single index or a set of indexes.

## Shard sizing
- Shard sizing depends on how much data the index holds (including future requirements) and home much heap memory is allocated to a node.
- An individual shard should be no more than 50GB and a node should only host up to 20 shards per 1GB of heap memory.

## Re-indexing


## Routing function
- Elasticsearch uses a routing function to distribute a given document to a shard when indexing: `shard_number = hash(id) % number_of_shards`.
- The same routing function is used to find the shard to which a document belongs.
- The routing function depends on the number of shards so this cannot be changed once an index is created.

## Clusters
- A cluster is a collection of nodes.
- A cluster can be in one of the following states:
  - Red: The shards have not yet been assigned, so not all data is available for querying.
  - Yellow: Replicas are not yet assigned, but all the shards have been assigned.
  - Green: All shards and replicas are assigned.
- Shards are readjusted across nodes in the case of node failure.

## Nodes
- A node is an instance of Elasticsearch.

## Node roles


## Scaling
- Elasticsearch supports both vertical and horizontal scaling.
- Data will be distributed to new nodes as the join the cluster.

## Inverted index
- A full-text field search query is tokenized and normalized using the same analyzers associated with that field and the resulting tokens are matched in the inverted index.

## Queries
- Data is retrieved via search and analytical queries.


