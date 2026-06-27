# Pagination

SPARK APIs provide two pagination strategies: cursor-based pagination for real-time, consistent results, and offset-based pagination for simple use cases.

## Cursor-Based Pagination (Default)

Cursor-based pagination is the default and recommended approach. Each page response includes a cursor pointing to the last item in the current page. Clients include this cursor in the next request to retrieve the next page. Cursors are opaque strings encoded as base64 and contain no semantic information. This approach handles real-time data correctly because new items inserted at the beginning of the result set do not cause page drift or duplicate results.

Cursor-based requests use the after and first parameters. The after parameter specifies the cursor to fetch results after. The first parameter specifies the maximum number of items to return, with a maximum of 100. The response includes edges (array of items), pageInfo (hasNextPage, hasPreviousPage, startCursor, endCursor), and totalCount.

## Offset-Based Pagination

Offset-based pagination uses page and perPage parameters. The page parameter is one-indexed and defaults to 1. The perPage parameter defaults to 20 with a maximum of 100. The response includes a meta object with currentPage, perPage, totalCount, and totalPages. Offset pagination is available for endpoints where data does not change frequently and consistent ordering is acceptable.

## Sorting

Sorting is specified through the sort parameter with the format {field}:{direction}. Multiple sort fields are comma-separated. Default sorting is by created_at:desc. Custom sort fields are documented per endpoint.
