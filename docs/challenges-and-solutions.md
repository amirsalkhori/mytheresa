# Challenges and Solutions

## Overview
This document outlines the challenges faced during the implementation of the MyTheresa Promotions Test and the corresponding solutions.

---

## Challenges

### 1. **Efficient Pagination**
- **Problem:** Handling large datasets with traditional offset-based pagination causes performance degradation as the offset increases.
- **Solution:** Implemented cursor-based pagination to avoid offset queries. Cursors use encoded product IDs to fetch subsequent pages efficiently. More in [implementation#Cursor-Based Pagination](implementations.md)

### 2. **Dynamic Discount Application**
- **Problem:** Applying discounts dynamically for `sku` and `category` with priority rules required frequent database lookups.
- **Solution:** Cached discount data in Redis, with fallback to the database on cache misses. Redis keys were structured for efficient lookups.

### 3. **Scalability**
- **Problem:** Ensuring the system can handle datasets exceeding 20,000 products while maintaining low response times.
- **Solution:**
  - Normalized database schema to avoid redundant data and improve query performance.
  - Indexed key fields (`sku`, `category_id`, and `price`) for faster filtering and sorting.
- **Performance analysis Results:** Available in [Performance](performance.md) section.

### 4. **Cache Invalidation**
- **Problem:** Maintaining cache consistency when products or discounts are updated.
- **Solution:** Applied TTLs to Redis entries to ensure eventual consistency. Future work could involve more precise invalidation mechanisms.
---

## Lessons Learned
- Caching is a powerful tool, but managing cache consistency is equally critical.
- Cursor-based pagination is an effective solution for large datasets when total count information is not required.


