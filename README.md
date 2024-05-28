## Overview
The service allows sending banners to users based on the requested feature and user tag. It is also possible to manage banners, features, and tags related to them.
## Definitions and relationships
**Banner** is a document describing an element of the user interface. Technically, a banner is a JSON document with an undefined structure.
**Tag** is an entity used to denote a group of users; it is represented by a number (tag ID).
**Feature** is a domain or functionality; it is represented by a number (feature ID).
1. One banner can be associated with only one feature and multiple tags.
2. However, one tag, like one feature, can belong to different banners simultaneously.
3. A feature and a tag uniquely identify a banner.
## Access, Management and Versioning
1. Two types of tokens are used for access authorization: user and admin. Banners can be retrieved using either a user token or an admin token, but all other actions can only be performed with an admin token.
2. Banners can be temporarily disabled. If a banner is disabled, regular users cannot access it, but admins still have access to it.
3. Sometimes it is necessary to revert to one of the three previous versions of a banner due to errors in logic, text, etc. It is possible to view the existing versions of the banner and choose the appropriate one.
4. Banners can be deleted by feature or tag.
### Detailed descriptions of the API methods can be found in the `openapi.yaml` file.