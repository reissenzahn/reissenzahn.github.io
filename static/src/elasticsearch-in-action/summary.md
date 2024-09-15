# Chapter 1
- Structured data is organized, has a definite shape and format and fits predefined data type patterns. It follows a defined schema and is easily searchable. Queries against structured data return results in exact matches.
- Unstructured data is unorganized and follows no schema or format. It has no predefined structure. Full-text (unstructured) queries try to find results that are relevant to the query.
- Relevancy refers to the degree to which search results match the query. A relevancy scores is a positive number given to each result that indicates how closely the result matches the query.
- Elasticsearch can index and search both structured and unstructured data in the same index.
- The Elastic Stack consists of:
  - Elasticsearch: A search and analytics engine built on Lucene; communication takes place over REST APIs using JSON.
  - Beats: Single-purpose data shippers that load data from various external systems and pump it into Elasticsearch.
  - Logstash: An ETL engine that extracts data from multiple sources, processes it and loads it into Elasticsearch.
  - Kibana: Web console that provides functionality for query execution, dashboards and visualizations.
- Elasticsearch is eventually consistent and optimized for high-read workloads.

# Chapter 2
- Elasticsearch is a document data store and it expects documents to be presented in JSON format.
- Communication with Elasticsearch takes place via JSON-based REST APIs.
- Documents are identified by an ID and are organized into indexes.
- An index is a bucket for collection documents. An index is expected to be dedicated to one and only one type.
- The _doc endpoint is a constant part of the path that is associated with the operation being performed.
- The document ID can be used to retrieve the document.
- Elasticsearch is schema-less.

```json
PUT books/_doc/1
{
  "title": "Effective Java",
  "author": "Joshua Bloch",
  "release_date": "2001-06-01",
  "amazon_rating": 4.7,
  "best_seller": true,
  "prices": {
    "usd": 9.95,
    "gbp": 7.95,
    "eur": 8.95
  }
}

PUT books/_doc/2 
{
  "title":"Core Java Volume I - Fundamentals",
  "author":"Cay S. Horstmann",
  "release_date":"2018-08-27",
  "amazon_rating":4.8,
  "best_seller":true,
  "prices": {
    "usd":19.95,
    "gbp":17.95,
    "eur":18.95
  }
}

PUT books/_doc/3 
{
  "title":"Java: A Beginner’s Guide",
  "author":"Herbert Schildt",
  "release_date":"2018-11-20",
  "amazon_rating":4.2,
  "best_seller":true,
  "prices": {
    "usd":19.99,
    "gbp":19.99,
    "eur":19.99
  }
}
```

```json
// count the total number of documents in an index
GET books/_count

// count the total number of documents in multiple indices
GET books,movies/_count

// count the total number of documents in all indices (including system and hidden indices)
GET _count
```

- Operations to retrieve documents require specifying the document ID.
- Elasticsearch provides Query DSL for writing queries in JSON.

```json
// retrieve a single document by ID
GET books/_doc/1

// exclude metadata in response
GET books/_source/1

// retrieve multiple documents given a set of IDs
GET books/_search
{
  "query": {
    "ids": {
      "values": [1,2,3]
    }
  }
}

// exclude source data in response
GET books/_search
{
  "_source": false, 
  "query": {
    "ids": {
      "values": [1,2,3]
    }
  }
}
```

```json
// retrieve all documents in index (short form)
GET books/_search

// retrieve all documents in index
GET books/_search
{
  "query": {
    "match_all": {}
  }
}

// search for words across full text
GET books/_search
{
  "query": {
    "match": {
      "author": "Joshua"
    }
  }
}

// search for documents where field contains a term that begins with a prefix
{
  "query": {
    "prefix": {
      "author": "josh"
    }
  }
}

// search for documents where field matches "Joshua" or "Schildt" (the OR operator is used implicitly)
GET books/_search
{
  "query": {
   "match": {
     "author": "Joshua Schildt" 
   }
  }
}

// search for documents where field matches "Joshua" and "Schildt"
GET books/_search
{
  "query": {
    "match": {
      "author": {
        "query": "Joshua Schildt",
        "operator": "AND",
      }
    }
  }
}
```

- The `_bulk` API allows for indexing documents simultaneously.
- Every two lines corresponds to one document. The first line is metadata about the record (operation, document ID) and the second line is the actual document.

```json
POST _bulk
{"index":{"_index":"books","_id":"1"}}
{"title": "Core Java Volume I – Fundamentals","author": "Cay S. Horstmann","edition": 11, "synopsis": "Java reference book that offers a detailed explanation of various features of Core Java, including exception handling, interfaces, and lambda expressions. Significant highlights of the book include simple language, conciseness, and detailed examples.","amazon_rating": 4.6,"release_date": "2018-08-27","tags": ["Programming Languages, Java Programming"]}
{"index":{"_index":"books","_id":"2"}}
{"title": "Effective Java","author": "Joshua Bloch", "edition": 3,"synopsis": "A must-have book for every Java programmer and Java aspirant, Effective Java makes up for an excellent complementary read with other Java books or learning material. The book offers 78 best practices to follow for making the code better.", "amazon_rating": 4.7, "release_date": "2017-12-27", "tags": ["Object Oriented Software Design"]}

// generate a random ID for each document
{"index":{}})
```

```json
// search across multiple fields
GET books/_search
{
  "query": {
    "multi_match": { 
      "query": "Java", 
      "fields": ["title","synopsis"] 
    }
  }
}

// give higher relevance to certain fields by providing a boost factor
GET books/_search
{
  "query": {
    "multi_match": { 
      "query": "Java",
      "fields": ["title^3","synopsis"]
    }
  }
}
```

- match_phrase expects a full phrase; the slop parameter indicates how many words the phrase is missing when searching

```json
// match a sequence of words exactly in a given order
GET books/_search
{
  "query": {
    "match_phrase": {  
      "synopsis": "must-have book for every Java programmer"  
    }
  }
}

// query can expect up to two words to be either missing or not in order when the query is executed
GET books/_search
{
  "query": {
    "match_phrase": {
      "synopsis": {
      "query": "must-have book every Java programmer",
      "slop": 2
      }
    }
  }
}
```

```json
// highlight specified fields
GET books/_search
{
  "query": {
    "match_phrase": {
      "synopsis": "must-have book for every Java programmer"
    }
  },
  "highlight": {
    "fields": {
      "synopsis": {}
    }
  }
}
```

- If fuzziness is set to 1, one spelling mistake (one letter misplaced, omitted, or added) can be forgiven.

```json
GET books/_search
{
  "query": {
    "match": {
      "tags": {
        "query": "Komputer",
        "fuzziness": 1 
      }
    }
  }
}
```

- In addition to full-text queries, Elasticsearch supports queries for searching structured data: term-level queries.
- Numbers, dates, ranges, IP addresses, and so on belong to the structured text category.
- When indexing a document for the first time, Elasticsearch deuces the schema by analyzing the values of the fields.
- Fields that are categorized as non-text fields are stored as they are.
- Term-level queries produce binary output: results are fetched if the query matches the criteria; otherwise, no results are sent.
- The range query can be applied to dates, numerals, and other attributes, making it a powerful companion when searching for range data.

```json
// a term query is used to fetch exact matches for a value provided in the search criteria.
GET books/_search
{
  "_source": ["title", "edition"],
  "query": {
    "term": {
      "edition": {
        "value": 3
      }
    }
  }
}

// a range query fetches results that match a range
GET books/_search
{
  "query": {
    "range": { 
      "amazon_rating": {
        "gte": 4.5, 
        "lte": 5 
      }
    }
  }
}
```

- Compound queries provide a mechanism to combine individual (leaf) queries.
- The compound queries are: bool, constant_score, function_score, boosting, dis_max.
- A Boolean query is used to combine other queries based on Boolean conditions. The search is composed of a set of four clauses:
  - must: The must clause means the search criteria in a query must match the documents. A positive match bumps up the relevancy score. We build a must clause with as many leaf queries as possible.
  - must_not: In a must_not clause, the criteria must not match the documents. This clause does not contribute to the score (it is run in a filter execution context).
  - should: It is not mandatory for a criterion defined in the should clause to match. However, if it matches, the relevancy score is bumped up.
  - filter: In the filter clause, the criteria must match the documents, similar to the must clause. The only difference is that the score is irrelevant in the filter clause (it is run in a filter execution context).
- The should clause behaves like an OR operator. That is, if the search words are matched against the should query, the relevancy score is bumped up. If the words do not match, the query does not fail; the clause is ignored. The should clause is more about increasing the relevancy score than affecting the results.

```json
GET books/_search
{
  "query": {
    "bool": {
      "must": [
        {"match": {"author": "Joshua Bloch"}},
        {"match_phrase": { "synopsis": "best Java programming books" }}
      ],
      "must_not": [
        {"range": {"amazon_rating": {"lt": 4.7 }}}
      ],
      "should": [
        {"match": {"tags": "Software"}}
      ],
      "filter": [
        {"range": {"release_date": {"gte":"2015-01-01"}}},
        {"term": {"edition": 3}}
      ]
    }
  }
}
```

- Aggregations fall into three categories:
  - Metric: Simple aggregations like sum, min, max and avg. Provide an aggregate value across a set of document data.
  - Bucket: Collect data into buckets segregated by intervals such as days, age groups, etc.
  - Pipeline: Work on the output from other aggregations.

```json
GET covid/_search
{
  "aggs": {
    "critical_patients": {
      "sum": {
        "field": "critical"
      }
    }
  }
}

// exclude source documents from response
GET books/_search
{
  "size": 0,
  "aggs": { 
    "avg_rating": {
      "avg": {
        "field": "amazon_rating" 
      }
    }
  }
}

// return basic statistical metrics
GET covid/_search
{
  "size": 0,
  "aggs": {
    "all_stats": {
      "stats": {
        "field": "deaths"
      }
    }
  }
}
```

```json
// categorize countries by the number of critical patients in buckets of 2500
GET covid/_search
{
  "size": 0,
  "aggs": {
    "critical_patients_as_histogram": {
      "histogram": {
        "field": "critical",
        "interval": 2500
      }
    }
  }
}

// range bucketing defines a set of buckets based on predefined ranges
GET covid/_search
{
  "size": 0,
  "aggs": {
    "range_countries": {
      "range": {
        "field": "deaths",
        "ranges": [
          {"to": 60000},
          {"from": 60000, "to": 70000},
          {"from": 70000, "to": 80000},
          {"from": 80000, "to": 120000},
        ]
      }
    }
  }
}
```

# Chapter 3
- An alias is an alternate name for a single index or a set of multiple indexes. The ideal way to search against multiple indexes is to create an alias  pointing  to  them.
- When  we  search  against  an  alias,  we  are  essentially searching against all the indexes managed by this alias. 
- Data streams accommodate time-series data: they let us hold data in multiple indexes but allow access as a single resource for search and analytical related queries.
- Each data stream has a set of indexes for each time point. These indexes are generated automatically and hidden.
- The data stream itself is nothing more than an alias for the time-series (rolling) hidden indexes behind the scenes. While search/read requests span all of the data stream’s backing hidden indexes, indexing requests are only directed to the new (current) index.
- Data streams are created using a matching indexing template. Templates are blueprints consisting of settings and configuration values used for creating resources like indexes. Indexes created from a template inherit the settings defined in the template.

- Shards hold data, created support data structures (like inverted indexes), manage queries and analyze data in Elasticsearch. They are instances of Lucene allocated to an index during index creation.
- During the process of indexing, a document travels through to the shard. Shards create immutable file segments to hold the document in a durable file system.
- Shards are distributed across the cluster for availability and failover.
- Once an index is in operation, shards cannot be relocated, as doing so would invalidate the existing data in an index. 
- As duplicate copies of shards, replica shards provide redundancy and high availability in an application. By serving read requests, replicas enable distributing the read load during peak times.
- Shards and their respective replicas are distributed across different nodes in the cluster. 


```json
GET _cluster/health
```


- Replicas will never be created on the same machine where their respective shards exist.
- When you boot up a node for the first time, Elasticsearch forms a new cluster, called a  single-node cluster.
- If we start another node in the same network, the newly instantiated node joins the existing cluster, provided the cluster.name property points to the single-node cluster.
- Adding more nodes to the cluster not only creates redundancy and makes the system fault tolerant but also brings huge performance benefits. When we add more nodes, we  create  more  room  for  replicas. 
- Usually, data is reindexed from an existing index into a new index. This new index  is  configured  with  a  new  number  of  shards,  taking  the  additional  new  nodes into consideration.

- For example, if we have a three-shards-and-15-replicas-per-shard strategy, with each shard sized at 50 GB, we must ensure that all 15 replicas have enough capacity not only to store documents on disk but also for heap memory: 
  - Shard memory: 3 ✕ 50 GB/shard = 150 GB 
  - Replica memory per shard: 15 ✕ 50 GB/replica = 750 GB/per shard
  - Total memory for both shards and replicas on a given node = 150 GB + 750 GB = 900 GB
  - Grand total for 20 nodes = 18 TB

- We can create a cluster of any number of nodes based on our data requirements.
- Bundling all sorts of data into a single cluster is not unusual, but it may not be a best practice. It
may be a better strategy to create multiple clusters for varied data shapes, with customized configurations for each cluster.
- Every  node  plays  multiple  roles,  from  being  a  coordinator  to  managing  data  to becoming a master.
  - Master node: Its primary responsibility is cluster management.
  - Data node: Responsible for document persistence and retrieval.
  - Ingest node: Responsible for the transformation of data via pipeline ingestion before indexing.
  - Machine learning: node Handles machine learning jobs and requests.
  - Transform node: Handles transformation requests.
  - Coordination node: This is the default role. It takes care of incoming client requests.

- If this master node crashes, the cluster elects one of the other nodes as the master, so the baton passes. Master nodes don’t participate in document CRUD operations, but the master node knows the location of the documents. 

CONFIGURING ROLES
Remember, we mentioned that the coordinator role is the default role provided to all
nodes. Although we set up four roles in the example (master, data, ingest, and ml),
this node still inherits a coordinator role. 
We can specifically assign just a coordinator role to a node by simply omitting any
node.roles  values.  In  the  following  snippet,  we  assign  the  node  nothing  but  the
coordinator role, meaning this node doesn’t participate in any activities other than
coordinating requests:
nodes.roles : [ ]  
There is a benefit to enabling nodes as dedicated coordinators: they perform as load
balancers,  working  through  requests  and  collating  the  result  sets.  However,  the  risk
outweighs the benefits if we enable many nodes as just coordinators. 
 In our earlier discussions, I mentioned that Elasticsearch stores analyzed full-text
fields in an advanced data structure called an inverted index. If there’s one data struc-
ture  that  any  search  engine  (not  just  Elasticsearch)  depends  on  heavily,  that’s  the
inverted  index.  It  is  time  to  examine  the  internal  workings  of  an  inverted  index,  to
solidify our understanding of the text analysis process, storage, and retrieval. 




- Elasticsearch uses a data structure called an inverted index for each full-text field during the indexing phase.
- Elasticsearch consults the inverted index for search words and document associations. When the engine finds document identifiers for these search words, it returns the full document(s).
- For each document that consists of full-text fields, the server creates a respective inverted index.


The analysis process is a complex function carried out by an analyzer
module.  The  analyzer  module  is  further  composed  of  character  filters,  a  tokenizer,
and token filters. When the first document is indexed, as in the greeting field (a text
field),  an  inverted  index  is  created.  Every  full-text  field  is  backed  up  by  an  inverted
index. The value of the greeting “Hello, World” is analyzed, tokenized, and normal-
ized into two words—hello and world—by the end of the process. But there are a few
steps in between. 
 Let’s look at the overall process (figure 3.18). The input line <h2>Hello WORLD</
h2> is stripped of unwanted characters such as HTML markup. The cleaned-up data is
split  into  tokens  (most  likely  individual  words)  based  on  whitespace,  thus  forming
Hello WORLD
[Hello] [WORLD]
[hello] [world]
Tokenizes the words
 (whitespace 
tokenizer)
Lowercases the 
tokens
(lowercase 
token filter)
<h2>Hello WORLD</h2>
Strips the HTML markup
(html_strip 
character filter)
Figure 3.18 Text analysis procedure where Elasticsearch processes text
88 CHAPTER 3 Architecture
Hello WORLD. Finally, token filters are applied so that the sentence can be transformed
into  tokens:  [hello]  [world].  By  default,  Elasticsearch  uses  a  standard  analyzer  to
lowercase the tokens, as is the case here. Note that the punctuation (comma) was also
removed during this process.
  After  these  steps,  an  inverted  index  is  created  for  this  field.  It  is  predominantly
used in full-text search functionalities. In essence, it is a hash map with words as keys
pointing  to  documents  where  these  words  are  present.  It  consists  of  a  set  of  unique
words  and  the  frequency  with  which  those  words  occur  across  all  documents  in  the
index. 
 Let’s revisit our example. Because document 1 (“Hello WORLD”) was indexed and
analyzed,  an  inverted  index  was  created  with  the  tokens  (individual  words)  and  the
documents they occur in; see table 3.3.
The words hello and  world are added to the inverted index along with the document
IDs  where  these  words  are  found  (of  course,  document  ID  1).  We  also  note the  fre-
quency of the words across all documents in this inverted index. 
  When  the  second  document  (“Hello,  Mate”)  is  indexed,  the  data  structure  is 
updated (table 3.4).
When updating the inverted index for the word hello, the document ID of the second
document is appended, and the frequency of the word increases. All tokens from the
incoming  document  are  matched  with  keys  in  the  inverted  index  before  being
appended to the data structure (like mate in this example) as a new record if the token
is seen for the first time.
  Now  that  the  inverted  index  has  been  created  from  these  documents,  when  a
search  for  hello  comes  along,  Elasticsearch  will  consult  this  inverted  index  first.  The
inverted index points out that the word hello is present in document IDs 1 and 2, so
the relevant documents will be fetched and returned to the client. 
Table 3.3 Tokenized words for the “Hello, World” document
Word Frequency Document ID
hello 1 1
world 1 1
Table 3.4 Tokenized words for the “Hello, Mate” document
Word Frequency Document ID
hello 2 1,2
world 1 1
mate 1 2
893.


While an inverted index is optimized for faster information retrieval, it adds com-
plexity  to  analysis  and  requires  more  space.  The  inverted  index  grows  as  indexing
activity increases, thus consuming computing resources and heap  space.

4 Relevancy
The inverted index also helps deduce relevancy scores; it provides the frequency of
terms, which is one of the ingredients in calculating the relevancy score. We’ve used
the  term  relevancy  for  a  while,  and  it  is  time  to  understand  what  it  is  and  what  algo-
rithms a search engine like Elasticsearch uses to fetch relevant results for the user. We
discuss concepts related to relevancy in the next section.
occur more frequently in a specific document and less frequently in the entire
collection should be considered more relevant as per TF-IDF algorithm.
 BM25 (Best Matching 25) algorithm. The BM25 relevancy algorithm is an improve-
ment over the TF-IDF algorithm. It introduces two important modifications to
the  base  algorithm  in  that  it  uses  a  nonlinear  function  for  term  frequency  to
prevent  highly  repetitive  terms  from  receiving  excessively  high  scores.  It  also
employs  a  document  length  normalization  factor  to  counter  the  bias  towards
longer documents. In the TF-IDF algorithm, longer documents are more likely
to have higher term frequencies, which means they may receive unduly credit
with higher scores. BM25’s job is to avoid such bias.
So,  both  BM25  and  TF-IDF  are  relevancy  algorithms  used  in  Elasticsearch.  BM25  is
considered  an  improvement  over  TF-IDF  due  to  its  term  frequency  saturation  and
document length normalization features. As a result, the BM25 is expected to be more
accurate and return relevant search results.
  Elasticsearch  provides  a  module  called  similarity  that  lets  us  apply  the  most
appropriate algorithm if the default isn’t suited to our requirements. 
 Similarity algorithms are applied per field by using mapping APIs. Because Elastic-
search  is  flexible,  it  also  allows  customized  algorithms  based  on  our  requirements.
(This is an advanced feature, so unfortunately we do not discuss it much in this book.)
Table 3.5 lists the algorithms available out of the box.
Table 3.5 Elasticsearch’s similarity algorithms
Similarity algorithm Type Description
Okapi BM25 
(default)
BM25 An enhanced TF/IDF algorithm that considers field length in 
addition to term and document frequencies
Divergence from 
Randomness (DFR)
DFR Uses the DFR framework developed by its authors, Amati and 
Rijsbergen, which aims to improve search relevance by measur-
ing the divergence between the actual term and an expected 
random distribution. Terms that occur more often in relevant 
documents than in a random distribution are assigned higher 
weights when ranking search results.
Divergence from 
Independence (DFI) 
DFI A specific model of the DFR family that measures the diver-
gence of the actual term frequency distribution from an inde-
pendent distribution. DFI aims to assign higher scores to 
documents by comparing the observed term frequencies with 
those expected in a random, uncorrelated term frequency.
LM Dirichlet LMDirichlet Calculates the relevance of documents based on the probability 
of generating the query terms from the document’s language 
model
LM Jelinek-Mercer LMJelinek-
Mercer
Provides improved search result relevance compared to models 
that do not account for data sparsity
Manua Creates a manual script
Boolean similarity boolean Does not consider ranking factors unless the query criteria are 
satisfied
92 CHAPTER 3 Architecture
In the next section, we briefly review the BM25 algorithm, which is a next-generation
enhanced TF/IDF algorithm.
THE OKAPI BM25 ALGORITHM
Three main factors are involved in associating a relevancy score with the results: term
frequency (TF), inverse document frequency (IDF), and field-length norm. Let’s look
at these factors briefly and learn how they affect relevancy.
 Term frequency represents the number of times the search word appears in the cur-
rent document’s field. If we search for a word in a title field, the number of times the
word  appears  is  denoted  by  the  term  frequency  variable.  The  higher  the  frequency,
the higher the score. 
 Say, for example, that we are searching for the word Java in a title field across three
documents. When indexing, we created the inverted index with similar information:
the word, the number of times that word appears in that field (in a document), and
the document IDs. We can create a table with this data, as table 3.6 shows.
Java appears three times in the document with ID 25 and one time in the other two
documents.  Because  the  search  word  appears  more  often  in  the  first  document  (ID
25), it is logical to consider that document our favorite. Remember, the higher the fre-
quency, the greater the relevance.
 While this number seems to be a pretty good indication of the most relevant docu-
ment in our search result, it is often not enough. Another factor, inverse document fre-
quency, when combined with TF, produces improved scores. 
 The number of times the search word appears across the whole set of documents
(i.e., across the whole index) is the document frequency. If the document frequency of a
word is high, we can deduce that the search word is common across the index. If the
word  appears  multiple  times  across  all  the  documents  in  an  index,  it  is  a  common
term and, accordingly, not that relevant. 
 Words that appear often are not significant: words like a, an, the, it, and so forth are
common in natural language and hence can be ignored. The inverse of the document
frequency  (inverse  document  frequency)  provides  a  higher  significance  for  uncommon
words across the whole index. Hence, the higher the document frequency, the lower the rele-
vancy, and vice versa. Table 3.7 shows the relationship between word frequency and rel-
evance.
Table 3.6 Term frequency for a search keyword
Title Frequency Doc ID
Mastering Java: 
Learning Core Java and Enterprise Java With Examples
3 25 
Effective Java 1 13 
Head First Java 1 39 
933.4 Relevancy
 
Until version 5.0, Elasticsearch used the TF-IDF similarity function to calculate scores
and rank results. The TF-IDF function was deprecated in favor of the BM25 function.
The TF-IDF algorithm didn’t consider the field’s length, which skewed the relevancy
scores. For example, which of these documents do you think is more relevant to the
search criteria?
 A field with 100 words, including 5 occurrences of a search word 
 A field with 10 words, including 3 occurrences of the search word
Logically,  it may  be  obvious  that  the  second  document  is  the  most  relevant  as it  has
more search words in a shorter length. Elasticsearch improved its similarity algorithms
by enhancing TF-IDF with an additional parameter: field length. 
 The field-length norm provides a score based on the length of the field: the search
word occurring multiple times in a short field is more relevant. For example, the word
Java  appearing  once  over  a  long  synopsis  may  not  indicate  a  useful  result.  On  the
Table 3.7 Relationship between word frequency and relevance
Word frequency Relevancy
Higher term frequency Higher relevancy
Higher document frequency Lower relevancy
Stop words
Words such as the, a, it, an, but, if, for, and and, are called stop words and can be
removed by using a stop filter plugin. The default standard analyzer doesn’t have the
stopwords parameter enabled (the stopwords filter is set to _none_ by default), so
these  words  are  analyzed.  However,  if  our  requirement  is  to  ignore  these  words,
we  can  enable  the  stop  words  filter  by  adding  the  parameter  stopwords  set  to
_english_, as shown here:
PUT index_with_stopwords
{
  "settings": {
    "analysis": {
      "analyzer": {
        "standard_with_stopwords_enabled": {
          "type": "standard",
          "stopwords": "_english_"
        }
      }
    }
  }
}
We learn about customizing analyzers in chapter 7.
94 CHAPTER 3 Architecture
other hand, as shown in table 3.8, the same word appearing twice or more in the title
field (with fewer words) says that the book is about the Java programming language. 
In  most  cases,  the  BM25  algorithm  is  adequate.  However,  if  we  need  to  swap  BM25
with another algorithm, we can do so by configuring it using the indexing APIs. Let’s
go over the mechanics of configuring the algorithm as needed.
CONFIGURING SIMILARITY ALGORITHMS
Elasticsearch  allows  us  to  plug  in  other  similarity  algorithms  if  the  default  BM25
doesn’t suit our requirements. Two similarity algorithms are provided out of the box
without further customization: BM25 and boolean. We can set the similarity algorithm
for individual fields when we create the schema definitions using index settings APIs,
as shown in figure 3.20.
NOTE Working with similarity algorithms is an advanced topic. While I advise
you to read this section, you can skip it and revisit it when you wish to know
more.  
Figure 3.20 Setting fields with different similarity functions
In the figure, an index index_with_different_similarities is being developed with
a schema that has two fields: title and author. The important point is the specifica-
Table 3.8 Comparing different fields to gather similarity
Word Field length Frequency Relevant?
Java Synopsis field with a length of 100 words 1 No
Java Title field with a length of 5 words 2 Yes
PUT index_with_different_similarities
{
 "mappings": {
   "properties": {
     "title":{
       "type": "text",
       "similarity": "BM25"
     },
     "author":{
       "type": "text",
       "similarity": "boolean"
     }
   }
 }
}
Creates an index with two
fields 
The title field is defined with a 
BM25 (default) similarity explicitly. 
The author field is defined
with a boolean similarity 
function. 
953.4 Relevancy
tion of two different algorithms attached  to  these  two fields  independently: title is
associated with the BM25 algorithm, while text is set with boolean. 
  Each  similarity  function  has  additional  parameters,  and  we  can  alter  them  to
reflect  precise  search  results.  For  example,  although  the  BM25  function  is  set  by
default  with  the  optimal  parameters,  we  can  easily  modify  the  function  using  the
index settings API. We can change two parameters in BM25 if we need to:  k1 and b,
described in table 3.9.
Let’s  look  at  an  example.  Figure  3.21  shows  an  index  with  a  custom  similarity  func-
tion, where the core BM25 function is amended with our own settings for k1 and b.
Figure 3.21 Setting custom parameters on the BM25 similarity function
Here, we are creating a custom similarity type—a tweaked version of BM25, which can
be reused elsewhere. (It’s more like a data type function, predefined and ready to be
attached to attributes.) Once this similarity function is created, we can use it when set-
ting up a field, as shown in figure 3.22.
Table 3.9 Available BM25 similarity function parameters
Property Default value Description
k1 1.2 Nonlinear term frequency saturation variable
b 0.75 TF normalization factor based on the document’s length
PUT my_bm25_index
{
 "settings": {
   "index":{
     "similarity":{
       "custom_BM25":{
         "type":"BM25",
         "k1":"1.1",
         "b":"0.85"
       }
     }
   }
 }
}
Configures an index with specific
BM25 parameters 
Sets the k1 and b values 
based on our requirements 
Creates a custom
similarity function with a
modified BM25 algorithm  
96 CHAPTER 3 Architecture
We  create  a  mapping  definition,  assigning  our  custom  similarity  function  (custom_
BM25) to a synopsis field. When ranking results based on this field, Elasticsearch con-
siders the provided custom similarity function to apply the scores.
You may wonder how Elasticsearch can retrieve documents in a fraction of a second.
How does it know where in these multiple shards the document exists? The key is the
routing algorithm, discussed next.






- Reindexing effectively creates a new index with appropriate settings and copies the data from the old
index to the new index. 
- Remember, the routing function is a function of the number of primary shards, not replicas. If we need to change the shard number, we must close the indexes (closed indexes block all read and write operations), change the shard number, and reopen the indexes.
- Alternatively, we can create a new index with a new set of shards and reindex the data from the old index to
the new index.




Summary
 Data is retrieved or searched via the search APIs (along with the document APIs
for single document retrieval).
 Incoming  data  must  be  wrapped  up  in  a  JSON  document.  Because  the  JSON
document is the fundamental data-holding  entity, it is persisted to shards and
replicas. 
 Shards and replicas are Apache Lucene instances whose responsibility is to per-
sist, retrieve, and distribute documents. 
 When we start up the Elasticsearch application, it boots up as a one-node, sin-
gle-cluster  application.  Adding  nodes  expands  the  cluster,  making  it  a  multi-
node cluster.
 For  faster  information  retrieval  and  data  persistence,  Elasticsearch  provides
advanced data structures like inverted indexes for structural data (such as tex-
tual  information)  and  BKD  trees  for  nonstructural  data  (such  as  dates  and
numbers).
 Relevancy  is  a  positive  floating-point  score  attached  to  retrieved  document
results. It defines how well the document matches the search criteria. 


# Chapter 4

- Data is modeled and indexed as JSON documents with each document consisting of a number of fields and every field containing a certain type of data.


It is mandatory to have the table schema defined and developed in a database before
retrieving or persisting data. But we can prime Elasticsearch with documents without
defining a schema for our data model. This schema-free feature helps developers get
up and running with the system from day one. However, best practice is to develop a
schema up front rather than letting Elasticsearch define it for us, unless our require-
ments do not need one. 
  Elasticsearch  expects  us  to  provide  clues  about  how  it  should  treat  a  field  when
indexing data. These clues are either provided by us in the form of a schema defini-
tion while creating the index or implicitly derived by the engine if we allow it to do so.
This process of creating the schema definition is called mapping.
 Mapping allows Elasticsearch to understand the shape of the data so it can apply a
set of predefined rules on fields before indexing them. Elasticsearch also consults the
manual of mapping rules to apply full-text rules on text fields. Structured fields (exact
values, like numbers or dates) have a separate set of instructions enabling them to be
part of aggregations and sorting and filtering functions, in addition to being available
for general searches.
 In this chapter, we set the context for using mapping schemas, explore the map-
ping process, and work with data types, looking at how to define them using the map-
ping  APIs.  Data  that  is  indexed  for  Elasticsearch  has  a  definite  shape  and  form.
Meticulous shaping of data lets Elasticsearch do a faultless analysis, providing the end
user with precise results. This chapter discusses the treatment of data in Elasticsearch
and how mapping schemas help us avoid hindrances and obtain accurate searches.
4.1 Overview of mapping
Mapping is a process of defining and developing a schema definition representing a
document’s data fields and their associated data types. Mapping tells the engine the
shape and form of the data that’s being indexed. Because Elasticsearch is a document-
oriented  database,  it  expects  a  single  mapping  definition  per  index.  Every  field  is
treated according to the mapping rule. For example, a string field is treated as a text
field, a number field is stored as an integer, a date field is indexed as a date to allow
for date-related operations, and so on. Accurate and error-free mapping allows Elastic-
search  to  analyze  data  faultlessly,  aiding  search-related  functionalities,  sorting,  filter-
ing, and aggregation.
100 m run or 400 m hurdles?
This chapter deals with dozens of hands-on examples around both core and advanced
data types. While I advise you to read about them in the given order, if you are just
starting with Elasticsearch and wish to focus on the beginner elements, you can skip
section 4.6 and revisit it when you are more confident and want to learn more. 
If all you want is a 100 m run rather than 400 m hurdles, read the chapter up to the
core data types (section 4.4) and then feel free to jump to the next chapter.
102 CHAPTER 4 Mapping
NOTE To simplify the coding exercises, I’ve created a ch04_mapping.txt file
under the kibana_scripts folder at the root of the book’s repository (http://
mng.bz/OxXo). Copy the contents of this file to Kibana as is. You can work
through the examples by executing the individual code snippets while follow-
ing the chapter’s content.
