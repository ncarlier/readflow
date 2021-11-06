ALTER TYPE article_status RENAME VALUE 'unread' TO 'inbox';
ALTER TYPE article_status ADD VALUE 'to_read' AFTER 'read';
