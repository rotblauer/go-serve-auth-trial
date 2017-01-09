I just wanted to follow up with the CSRF idea.

The majority of my experience with CSRF has been with pre-built solutions (ie Rails, golang libs, JWTs), so I spent some time today learning about it more in depth.

My takeaway is that my original `main.go` implements CSRF protection... but, by doing so with just a Cookie instead of Cookie-to-Header-Token[1]. My original implementation left exposed the possibility for a hijacker to forge a logout link, the only type of data/state-altering request made on behalf of the authenticated client. I addressed this by implementing a server-side check for a hidden token input field upon a logout POST, and creating a user-driven test in the DOM. 

It is also worth noting that my simple approach combines the idea of a Session with cross-origin forgery protection. At a higher scale the values accompanying authentication would become more commplex, accomodating and differentiating user and session data (you could store a key-value dictionary in the cookie, for example, and/or keep track of user states on the server side with something like Redis). 

[1]: https://en.wikipedia.org/wiki/Cross-site_request_forgery#Cookie-to-Header_Token 
[2]: http://stackoverflow.com/questions/20504846/why-is-it-common-to-put-csrf-prevention-tokens-in-cookies 

