#### 旧版本 v1.0.4
```go
	type MemoryGetRequestMemoryType string
	var (
	  MemoryGetRequestMemoryTypePerpetual        MemoryGetRequestMemoryType = "perpetual"
	  MemoryGetRequestMemoryTypeSummaryRetriever MemoryGetRequestMemoryType = "summary_retriever"
	  MemoryGetRequestMemoryTypeMessageWindow    MemoryGetRequestMemoryType = "message_window"```	
	)
````

#### 新版本 v1.0.6
```go
	type MemoryType string
	var (
	  MemoryTypePerpetual        MemoryType = "perpetual"
	  MemoryTypeSummaryRetriever MemoryType = "summary_retriever"
	  MemoryTypeMessageWindow    MemoryType = "message_window"
	)
```
