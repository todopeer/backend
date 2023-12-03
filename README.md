# TodoPeer Backend

The backend for TodoPeer App

See <https://api.todopeer.com/> for GraphQL APIs

## Concepts

For the basic model below, see the graphql page for the exact listing of fields

- Task: a tsak is something to be done. It has status (not_started / doing / done), description, due_date
- Event: 
    - an event is a small portion of a bigger task
    - an event belongs to a task; a task has many events
    - an event has a start & end time
    - it also has description, where one can describe what special about it


## UIs

There're currently 2 projects serving as UI:

- [CLI](https://github.com/todopeer/cli): CMD interface for defining the WebView
- [Frontend](https://github.com/todopeer/frontend): suppose to be used for both Web & Mobile