# 依赖注入

## 依赖注入的好处

简单来说，依赖注入是一种良好的编程规范，它有以下好处：
- 提供更规范、整洁的代码
- 形成更清晰的思路
- 执行 [Explicit Dependencies Principle](https://learn.microsoft.com/en-us/dotnet/architecture/modern-web-apps-azure/architectural-principles#explicit-dependencies)
- 促进单一职责原则 (Single Responsibility Principle)

更多依赖注入的基础知识，可以阅读以下材料：
- https://learn.microsoft.com/en-us/dotnet/architecture/modern-web-apps-azure/architectural-principles#dependency-inversion
- https://martinfowler.com/articles/injection.html


## 5分钟看懂：其他语言的依赖注入

目前广泛使用依赖注入的语言有java、C#。

在java的spring框架中，提供了依赖注入的方式：https://docs.spring.io/spring-framework/reference/core/beans/annotation-config/autowired.html

C#也提供了依赖注入的方式：https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection

简单来讲，依赖注入无非就是回答下列几个问题：
1. 作用域（Scope）：Service的作用域是什么？
2. 工厂方法（Factory method）：如何创建Service对象？创建Service对象前，需要依赖于其他什么Service？

搞清楚了这两个问题，我们就搞清楚了所有已知的、未知的语言的依赖注入。

我们先以 Asp.Net Core (C#) 的依赖注入为例。C#的依赖注入比java更简单、更好理解。

首先，我们需要建立一个Service。为了达成 Explicit Dependencies Principle，我们在它的构造函数中列出它依赖的其他Service。

```cs
public class ChatService
{
    private INetworkService networkService;

    public ChatService(INetworkService networkService)
    {
        this.networkService = networkService;
    }

    public Chat(string message)
    {
        var messageObject = new MessageObject(message, IPV6);
        networkService.GetConnection().SendMessage(messageObject);
    }
}
```

此时，`ChatService`的构造函数就回答了“工厂方法（Factory method）”的问题。还需要回答“作用域”的问题。但是，作用域放在Service上是不合适的，因为根据应用层场景的不同，Service的作用域也会不同。作用域需要在应用层定义，而非领域服务中定义。

在应用层定义作用域：

```cs
HostApplicationBuilder builder = Host.CreateApplicationBuilder(args);

builder.Services.AddSingleton<ChatService>();
```

上述代码向依赖注入框架提供了以下信息：
- 作用域：单例（Singleton）
- 服务：ChatService
- 工厂方法（Factory method）：ChatService的构造函数

以上信息能够指导框架创建出一个 ChatService 了。当我们向框架索要ChatService对象时，框架会使用我们提供的工厂方法进行创建。工厂方法可能会依赖于其他Service对象的创建，所以框架会以同样的方式，先试图得到其他Service对象。

此时，框架会产生疑问：`INetworkService`从哪里来？所以，我们也需要提供 `INetworkService` 的相关信息：

```cs
builder.Services.AddSingleton<INetworkService, StableNetworkService>();
```

上述代码提供的信息：
- 作用域：单例
- 服务：INetworkService
- 工厂方法：StableNetworkService的构造函数

有时我们可能已经持有`StableNetworkService`对象，想要直接使用这个对象，而不是重新调用工厂方法。例如，可能`StableNetworkService`的sdk以静态变量的方式提供了一个`sdk.DefaultStableNetworkService`对象。此时，我们可以这样做：

```cs
builder.Services.AddSingleton<INetworkService, StableNetworkService>(() => sdk.DefaultStableNetworkService);
```

上述代码提供的信息：
- 作用域：单例
- 服务：INetworkService
- 工厂方法：一个lambda函数，会直接返回 `sdk.DefaultStableNetworkService`

理解了上述内容，就理解了C#的依赖注入。是不是很简单？

让我们回过头来看java（Spring）的依赖注入。

构造函数：

```java
// https://docs.spring.io/spring-framework/reference/core/beans/annotation-config/autowired.html

@Component
@Scope("singleton")
public class MovieRecommender {

	private final CustomerPreferenceDao customerPreferenceDao;

	@Autowired
	public MovieRecommender(CustomerPreferenceDao customerPreferenceDao) {
		this.customerPreferenceDao = customerPreferenceDao;
	}
}
```

以上代码可以和C#的 `builder.Services.AddSingleton<MovieRecommender>()` 对应。
- `AddSingleton`: `@Component` `@Scope("singleton")`
- `<MovieRecommender>`:  `@Autowired public MovieRecommender(...`

可以看出，东西能够对应上。但是，spring将应用层的东西（@Component, @Scope）放到了领域服务中，这不是标准的、规范的做法，会存在一定的风险。

另外，如果想要实现类似 `builder.Services.AddSingleton<INetworkService, StableNetworkService>(() => sdk.DefaultStableNetworkService)` 的功能，可以使用如下的代码：

```java
@Configuration
public class Config {

    @Bean
    @Scope("singleton")
    public INetworkService provideNetworkService() {
        return sdk.DefaultStableNetworkService;
    }
}

```

以上代码就非常接近 `builder.Services.AddSingleton<INetworkService, StableNetworkService>(() => sdk.DefaultStableNetworkService)` 了。


## golang的依赖注入思路

在 java 和 C# 中，Spring 框架和 Asp.Net Core 框架，它们都有依赖注入的功能。在golang中，我们也有依赖注入的框架，例如wire、dig。

当我们想要创建一个 Sevice对象时，这些框架会找到工厂方法（Factory Method），然后看工厂方法的依赖的 Service对象是否满足，如果未满足，就会尝试创建这些 Service对象。Service的依赖关系构成了一张有向无环图，而框架会帮我们维护这张图，并在我们需要 Service对象时，通过这张图为我们创建对象。

那么，假如没有依赖注入框架，就做不了依赖注入吗？答案是否定的。

我们假设有这样2个存在依赖关系的 Service: AlphaService, BetaService.

首先，我们模仿C#和java的形式，遵循 Explicit Dependencies Principle，依赖的关系都写在构造函数里面。在golang里面并没有构造函数这个东西，但是我们常常将 `New` + Service 名称的形式作为约定俗成的构造函数。

services/alpha_beta.go
```go
type AlphaService struct {
}

func NewAlphaService() *AlphaService {
    return &AlphaService{}
}


type BetaService struct {
    alphaService *AlphaService
}

func NewBetaService(alphaService *AlphaService) *BetaService {
    return &BetaService{alphaService: alphaService}
}
```



注意，我们需要将几乎所有Service性质的全局变量，都通过构造函数进行中转，来满足 Explicit Dependencies Principle。在少数情况下，某些非常基础的Service（例如日志）可以不被考虑外，其他Service都强烈建议采用依赖注入的模式。

做到了这一步，golang的依赖注入就成功了一半。接下来，我们需要考虑如何手动管理 Service 对象的创建。Service 对象的创建不是领域服务需要关心的事情，而是应用层需要关系的事情。所以我们将这部分代码放到 `application/graph.go` 里面：

```go
func getDependencies() map[string]any {
    var alphaService = services.NewAlphaService()
    var betaService = services.NewBetaService(alphaService)

    return &map[string]any{
        "alphaService": alphaService,
        "betaService": betaService,
    }
}

var DependencyMap = getDependencies()
```

使用：

main.go

```go
func main() {
    var betaService = application.DependencyMap["betaService"].(services.BetaService)

    // var app = ...
    app.router.Register(betaService)

    app.run()
}
```

单元测试使用：

```go
func TestAbc(t *testing.T) {
    var betaService = application.DependencyMap["betaService"].(services.BetaService)

    betaService.SomeMethod()
}
```

当然，使用一个map来储存依赖，可能会不太方便，所以我们可以直接使用全局变量：

application/graph.go

```go
var alphaService = services.NewAlphaService()
var betaService = services.NewBetaService(alphaService)

// 在main包和测试包访问，不能访问私有变量，所以通过公共方法提供访问

func GetSingletonAlphaService() *AlphaService {return alphaService}
func GetSingletonBetaService() *BetaService {return betaService}
```

main.go

```go
func main() {
    var betaService = application.GetSingletonAlphaService()

    // var app = ...
    app.router.Register(betaService)

    app.run()
}
```

这样就变得更简洁了。


## 简单依赖注入

“简单依赖注入”并非一个已有的专有名词。它指的是这样的风格：

service/alpha_service.go
```go
type AlphaService struct {
}

// Constructor of AlphaService
func NewAlphaService() *AlphaService {
	return &AlphaService{}
}

// wire

var singletonAlphaService *AlphaService = initSingletonAlphaService()

func GetSingletonAlphaService() *AlphaService {
	return singletonAlphaService
}

func initSingletonAlphaService() *AlphaService {
	return NewAlphaService()
}
```

service/beta_service.go
```go
type BetaService struct {
	alphaService *AlphaService
}

// Constructor of BetaService
func NewBetaService(
	alphaService *AlphaService,
) *BetaService {
	return &BetaService{alphaService: alphaService}
}

// wire

var singletonBetaService *BetaService = initSingletonBetaService()

func GetSingletonBetaService() *BetaService {
	return singletonBetaService
}

func initSingletonBetaService() *BetaService {
	return NewBetaService(GetSingletonAlphaService())
}
```

“简单依赖注入”将Service对象的创建放到了服务定义的同一package、同一文件中，稍微有些违背上文说的只有application层才关心如何创建 Service对象。但同时，简单依赖注入具备一些优点：

- 与标准的做法相比，在代码仓库不被其他仓库引用时，不会存在不便
- 学习成本低，容易模仿

“简单依赖注入”其实算是一种折中。与标准做法相比，它更简单，且不采用任何框架。与完全不使用依赖注入相比，简单依赖注入的改造成本低，并且带来的收益高。


## golang依赖注入框架

## RequestScope 的依赖
