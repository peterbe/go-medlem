medlem
=======

About the project
-----------------

**THIS PROJECT IS IN PROTOTYPAL STAGES**

This project aims to set up a server for HTTP interfacing with LDAP.
Two core feature areas:

1) Be able to ask simple questions in HTTP GET and get a simple JSON response.
   For example "Is example@example.org still an employee of Example inc?"

2) If changes to membership happens, notify, by HTTP POST those who are
concerned to know this.

About the name
--------------

"medlem" means "member" in Swedish. Calling it just "member" would be too
common that it doesn't make for a good project name.

Also, LDAP it very depends on LDAP to get its source of truth. But that
might change in some future. At that point, the questions are still about
"membership".
