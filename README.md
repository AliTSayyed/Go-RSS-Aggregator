# RSS Aggregator with GO

The purpose of this project is to build an RSS Aggregator using GO for the backend and PostgreSQL as the database. RSS stands for Really Simple Syndication and an RSS aggregator is an application that collects RSS feeds from multiple sources and displays them in one place. User have the ability to select what RSS feeds they would like to follow and the aggregator will pull in those feeds and update the posts after a set amount of time so that the user's feed is up to date. 

An RSS feed is an online file that contains information about a website's published content. It can include details like the content's full text or summary, publication date, author, and link.

RSS feed's data is written in XML and we will automatically collect all the xml files from those feeds and save them into a database.
Then we can view the feeds and display them when / how we want. 

(Future additons)
1. Add functionality for different publishAt formats in scraper.go.
2. Add functionalty for different xml formats in rss.go.
3. Create custom front end to display posts. 