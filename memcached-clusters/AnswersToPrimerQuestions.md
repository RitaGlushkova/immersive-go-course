## Q: What is a chache hit rate and why is it important?

Chache hit rate describes the situation where content is served from cache and not from the original database. It is a percentage of a cachable data. It is important for a better performance of the app because cache uses RAM in comparison to disk memory, which is much slower.

## Q: What do we mean by ‘cold’ or ‘warm’ when we discuss chaches?

Is chache box is not filled with useful data it is ‘cold’. Every re## quest to this chache box will be a miss, so lots of resources will be used to fill up this chache. Warm chache is the one filled with useful for user data, so it can return ## quickly what the user is looking for.

## Q: Why do we use consistent hashing when we shard data?

A consistent hash function maps a key to a data partition and if a number of partition changes not all the keys has be remapped avoiding all chache go cold again.

## Q: When should we consider sharding a chache rather than replicating it?

When we are dealing with very large re## quests that can be bigger than any machine can handle. In this case multiple machines hosts a partition of a cached data. It allows us to go on a very big scale and server large load of data fast. It adds more complexity and work to a team to set it up.

## Q: Why do we need chache invalidation?

To avoid stale date returning to the user. If data in database has changed, and TTL is not expired, cache becomes invalid. So we need to invalidate cache and force new data to be fetched.

## Q: What is a TTL used for in chaching?

TTL - Time-To-Live responsible for how long to keep cache. Once it expires, data is removed from chache and must be fetched again from the database. It helps to keep information up to date and not overfill the cache box. It helps to make sure that data doesn’t remain in a stale state for to long.

## Q: Why should we cache in standalone cache service, as opposed to within our application?

If we do it within our application we have to make sure that a re## quest concerning a certain set of data must be routed to the same instance of our application (same server). We use sticky sessions for this matter but it limits our scalability (load barance can not manage load at its best efficiency) and it also if the server which holds our cache is down, the session will be lost.

## Q: Why are smaller data shards recommended?

It will help to find the correct data faster for a received request. If request hit a large database, it will have to search every row to find needed data. Another reason might be to reduce the risk of outages. If application relies on one or just a few large databases, then an outage may couse the whole application to stop working.

## Q:What are the challenges associated with recovering from a failed database replica?

Recovering leader. We need to designate a new machine as a leaders. It is usually most up-to-date replica but not always. Because leader is running on new machine all clients require reconfiguration. We also have to make sure that old leader doesn’t become online again as a leader (receiving updates from two leaders can be damage data in replicas).
If async replication is involved, all writes followers haven’t processed yet will be lost.

## Q: What kind of changes might you need to make to your application before moving to a sharded datastore?

Write sharding code and embed sharding logic in your application. But there are sharding schemes available.

## Q: What is a cell? What happens if a cell fails?

Cell is a set of servers that are very close together, and share the same regional availability
