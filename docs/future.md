# Future Improvements

## Overview

This document outlines potential future improvements to the MyTheresa Promotions Test implementation. These ideas were conceived during the development process but were not implemented due to time constraints. However, they represent opportunities to further enhance the system's performance, scalability, and maintainability.

---

## Redis-Based Filter Caching

### The Idea
To optimize the performance of filter queries, a Redis-based caching strategy could be implemented. This involves caching product IDs corresponding to specific filters, such as category or price range, directly in Redis.

### Implementation Details
1. **Caching Filter Results:**
   - For every filter query (e.g., `category=boots AND price<100000`), store the resulting product IDs in Redis.
   - Use Redis sets to map each filter to its associated product IDs. For example:
     - Key: `filter:category=boots:price<100000`
     - Value: `{product_id1, product_id2, product_id3, ...}`

2. **Cache Invalidation:**
   - Maintain a reverse mapping of products to the filters they belong to. For example:
     - Key: `product:123` 
     - Value: `{filter:category=boots:price<100000, filter:category=boots}`
   - When a product is updated (e.g., price changes or category changes), use this mapping to pinpoint the specific filters that need invalidation.
   - Remove the product ID from the relevant Redis sets or invalidate the filter entirely.
   - For example, if the price of `product_id3` changes, and it no longer fits `price<100000`, remove it from the associated Redis set.

3. **Cache Warming:**
   - Precompute and populate Redis with common filters during non-peak times or on application startup.
   - This reduces the cold cache penalty for frequently used filters.

4. **TTL (Time-to-Live):**
   - Apply TTL to cache entries to ensure they remain fresh and consistent with the database.
   - Example: Set a TTL of 1 hour for filter results and refresh them periodically.

### Benefits
- **Reduced Database Load:** Redis can serve filter results for frequent queries without requiring repeated database access.
- **Improved Query Performance:** Fetching results from Redis is significantly faster than executing database queries.
- **Scalability:** This approach scales well as the number of products and filters grows.
- **Targeted Invalidation:** Reverse mapping of products to filters ensures that only affected filters are invalidated, reducing unnecessary overhead.

### Challenges
- **Cache Invalidation Complexity:** Keeping Redis in sync with database updates requires careful planning and logic.
- **Memory Usage:** Large datasets or complex filters might increase Redis memory requirements.
- **Development Overhead:** Implementing and maintaining this system adds complexity to the codebase.

### Why It Was Not Implemented
While this approach could significantly improve performance, it was not implemented due to time constraints. The focus was on delivering a reliable MVP that meets the core requirements.

---

## Hybrid Pagination Strategy

### The Idea
A hybrid pagination approach combines traditional offset-based pagination for initial pages with cursor-based pagination for deeper pages. This approach provides the benefits of both methods:

- **Offset-based Pagination:** Used for the first few pages to give users a quick overview of available results.
- **Cursor-based Pagination:** Activated for deeper pages to maintain performance.

### Implementation Details
1. **Initial Caching for Offset-Based Pages:**
   - Fetch the first 100-200 records and cache them in Redis with their IDs.
   - Serve offset-based pages directly from Redis for immediate response.

2. **Switching to Cursor-Based Pagination:**
   - For pages beyond the pre-cached range, use cursor-based pagination.
   - Generate `next` and `prev` cursors for navigation.

3. **Integration with Frontend:**
   - Introduce infinite scrolling or a "load more" button to seamlessly transition users to cursor-based pagination.

### Benefits
- **User-Friendly Navigation:** Allows users to skip around early pages while still leveraging efficient pagination for large datasets.
- **Performance:** Maintains high performance for deep pagination.

### Challenges
- **Implementation Complexity:** Combining two pagination strategies requires additional logic and testing.
- **Frontend Coordination:** The UI must adapt dynamically based on the pagination mode.

### Why It Was Not Implemented
This approach requires additional infrastructure and logic, which could not be prioritized within the current timeline. It remains a promising area for future development.

---

## Elasticsearch Integration

### The Idea
Integrating Elasticsearch would enable advanced search capabilities, including full-text search, complex filtering, and near-instantaneous response times for queries.

### Implementation Details
1. **Indexing Data:**
   - Synchronize the `products` table with an Elasticsearch index.
   - Include fields such as `sku`, `name`, `category`, and `price`.

2. **Advanced Filtering:**
   - Use Elasticsearch's query DSL to perform complex queries, such as multi-field filtering and full-text search.

3. **Real-Time Updates:**
   - Implement a mechanism to update the Elasticsearch index whenever the database changes.
   - Use tools like RabbitMQ or a direct API integration for this purpose.

### Benefits
- **Enhanced Search:** Full-text search and more sophisticated filters compared to SQL queries.
- **High Scalability:** Designed to handle large datasets efficiently.
- **Improved Performance:** Faster response times for complex queries.

### Challenges
- **Infrastructure Overhead:** Elasticsearch requires additional setup and monitoring.
- **Data Synchronization:** Keeping Elasticsearch in sync with the primary database can be challenging.

### Why It Was Not Implemented
While Elasticsearch offers significant advantages, the current implementation focused on meeting core requirements without introducing additional infrastructure complexity.

---

## Conclusion

These future improvements demonstrate a clear path for scaling and optimizing the system. While they were not implemented due to time constraints, they reflect thoughtful consideration of potential challenges and solutions for a more robust application. These ideas could be revisited in future iterations as the system evolves.

