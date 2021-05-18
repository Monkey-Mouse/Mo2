# Comment system
> 2021/3/11 designed and implemented by chronos  

Basic: 二级评论结构

## API

- url: `/api/comment/{id}` GET
  - query: page, pagesize
  - get comments by article id

- url: `/api/comment/{id}` POST
  - publish a comment to an article by article id

- url: `/api/comment/sub/{id}` POST
  - publish a subcomment to a comment by comment id

- ~~url: `/api/comment/{id}` delete~~
  - delete a comment

- ~~url: `/api/comment/{id}/{subid}` delete~~
  - delete a sub comment

**KEEP IT SIMPLE**: we don't need delete. Frontend can change it's content to sth like "deleted" and upsert it  
===========================