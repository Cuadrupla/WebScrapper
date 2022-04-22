# WebScrapper
Challenge - This was a course challenge provided for us by the state of CS112 course.

The scrapper was built using the GoLang features and fast 💨 implementation.
The **authors** of the project are:
- Ispas Jany-Gabriel (group 142)
- Gheorghe Liviu-Ionut (group 144)
- Vrinceanu Radu-Tudor (group 142)

---

## How to run it?
The project has no dependency other than those provided under the GoLang standard libs so you can simply run it like this:
```
go run crawler.go --beginUrl="https://change.me" --depth=2
```

The results of the crawler will be shown in ```stdout``` like this:
```
[CRAWLER RESULTS]:
Used link: <beginUrl> (recursion depth: <depth>)
Found (total: x | unique: y | visited: z) URLS with crawler
```
