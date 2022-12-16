Porto Events
====

Collect Porto city events from popular resources and post theme to telegram chanel

Web server screenshot
![screenshot-web-server.jpg](screenshot-web-server.jpg)
Telegram channel post screenshot
![screenshot-telegram-post.jpg](screenshot-telegram-post.jpg)

## Features
- Web server: 
  - Showing collected events in list "New"
  - Init new events collect (add only new, not existed events)
  - Edit/Delete events
  - Move new event to "Publish" list
  - Init send "Publish" list to telegram
- Store events in DB (BoltDB)

#### Features under development
- Scheduler to collect new events automatically  
- Scheduler to publish events to telegram automatically  
- Add simple auth service for website
- Add config page (config schedulers, resources list enable/disable and so on)
- Add errors notification to the frontend

## Installation
- configs prepare
  - copy file config `configs/config-example.toml` an rename to `configs/config.toml`
  - add telegram token & channel id

## Store data explanation
At the beginning was idea store somewhere events. With ability see new events, edit and send to some list `to Publish`. From there get all events and send to telegram channel.

First try was not create own web server with DB and website for that.

Was chosen like a great idea store data in notes app NOTION. For easy manipulation of events data for client.
One note for configuration settings (how often collect events from resources, how often post to telegram and so on). 
One note for events list to edit, and add tag `#publish` if event is ready to publish to telegram  
Notion has APi so should be easy. But Notion has difficult note data structure. With different blocks. So not very easy to manipulate with data via api to save and after that to send to telegram.

**So** was chosen other solution: Create own web server and store data in light key/value Bolt DB without any additional installation requirements for the hosting


**Notion is under development stage. Not ready to use**

## Resources for get events
* [Porto.pt](https://www.porto.pt/en/events)
* [Agendaculturalporto.org](https://agendaculturalporto.org/agenda-maus-habitos-porto)
