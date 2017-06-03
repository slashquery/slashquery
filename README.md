# SlashQuery

API Gateway

The big picture:

```pre
   (slashquery API Gateway)

api.example.com/<path endpoint>
               |            (plugins)       (upstream)      (microservices)
               |                                            /home
               | /*     >--[middleware]---[LoadBalancer]--< /home
               |                                            /home
               |
               |                                            /foo
               | /foo/* >--[middleware]---[LoadBalancer]--< /foo
               |                                            /foo
               |
               |                                            /foo-version-2
 (versioning)  | /foo/* >--[middleware]---[LoadBalancer]--< /foo-version-2
               |                                            /foo-version-2
               |
               |                                            /bar
               ` /bar/* >--[middleware]---[LoadBalancer]--< /bar
                                                            /bar
```
