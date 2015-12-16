Uses SQLite full text search as described in https://sqlite.org/fts3.html.


The program expects a table with the following schema.

```sql
CREATE virtual TABLE data USING fts3(
	key,
	publisher,
	school,
	title,
	series,
	journal,
	author,
	number,
	month,
	volume,
	year,
	howpublished,
	organization,
	booktitle,
	institution
);
```

You can move data from a table with the following query.

```sql
INSERT INTO data select * from bibtex_entry;
```

Might need to recompile sqlite lib with `go install --tags "fts3"` if the module is missing.
