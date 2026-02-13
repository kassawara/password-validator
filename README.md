# Password Validator API

Uma aplicação web que valida senhas de acordo com critérios específicos de segurança. Desenvolvida em **Go** com arquitetura Clean Architecture e princípios SOLID, demonstrando boas práticas de engenharia de software.

## 📋 Requisitos de Validação

Uma senha é considerada válida quando possui:

- ✅ Nove ou mais caracteres (espaços não contam)
- ✅ Pelo menos 1 dígito (0-9)
- ✅ Pelo menos 1 letra minúscula (a-z)
- ✅ Pelo menos 1 letra maiúscula (A-Z)
- ✅ Pelo menos 1 caractere especial: `!@#$%^&*()-+`
- ✅ Sem caracteres repetidos (espaços são ignorados)

### Exemplos

```
IsValid("") // false  
IsValid("aa") // false  
IsValid("ab") // false  
IsValid("AAAbbbCc") // false  
IsValid("AbTp9!foo") // false  
IsValid("AbTp9!foA") // false
IsValid("AbTp9 fok") // false
IsValid("AbTp9!fok") // true
IsValid("  Abc def1!2  ") // true (espaços são ignorados)
```

---

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.25.0 ou superior
- Docker e Docker Compose (opcional)

### Instalação Local

1. Clone o repositório:
```bash
git clone <repository-url>
cd password-validator
```

2. Instale as dependências:
```bash
go mod download
```

3. Configure as variáveis de ambiente (opcional):
```bash
cp .env.example .env
# Edite .env conforme necessário
```

4. Execute a aplicação:
```bash
go run main.go
```

A API estará disponível em `http://localhost:8080`

### Com Docker Compose

```bash
docker-compose up --build
```

---

## 📡 API Endpoints

### Health Check
```http
GET /health HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

### Validar Senha
```http
POST /password/validate HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "password": "AbTp9!fok"
}
```

**Response (200 OK):**
```json
{
  "isValid": true
}
```

**Response (400 Bad Request):**
```json
{
  "errors": [
    {
      "field": "password",
      "message": "Must have at least 9 characters (excluding spaces)"
    }
  ]
}
```

---

## 🏗️ Arquitetura da Solução

A solução foi desenvolvida seguindo **Clean Architecture** com separação clara de responsabilidades:

```
password-validator/
├── adapter/                    # Adaptadores e Controllers
│   ├── controller/            # HTTP Controllers (Request/Response)
│   ├── handler/               # Tratamento de erros
│   ├── presenter/             # Formatação de output
│   ├── repository/            # Implementação do repositório
│   └── response/              # Estruturas de resposta HTTP
├── core/                       # Lógica de negócio
│   ├── domain/                # Entidades de domínio
│   │   └── password/          # Agregado Password
│   ├── usecase/               # Casos de uso
│   │   ├── input/             # DTOs de entrada
│   │   └── output/            # DTOs de saída
│   ├── repository/            # Interfaces de repositório
│   ├── errors/                # Erros customizados de domínio
│   └── utils/                 # Constantes e utilitários
├── infrastructure/            # Camada de infraestrutura
│   ├── config/                # Configurações da aplicação
│   ├── http/                  # Servidor HTTP
│   │   ├── server/            # Inicialização do servidor
│   │   ├── router/            # Definição de rotas
│   │   └── docs/              # Documentação Swagger
└── main.go                    # Ponto de entrada

```

---

## 🔑 Decisões de Arquitetura

### 1. **Clean Architecture**

A solução segue os princípios de Clean Architecture para garantir:
- **Independência de Frameworks**: A lógica de negócio não depende do Gin ou HTTP
- **Testabilidade**: Cada camada pode ser testada isoladamente
- **Manutenibilidade**: Mudanças em uma camada não afetam outras

### 2. **Domínio Anêmico com Validação**

O agregado `Password` concentra toda a lógica de validação:
```go
type Password struct {
    password string
    isValid  bool
}

func New(params ...PasswordParams) (*Password, error) {
    // Validação ocorre no construtor
    // Retorna erro se inválido
    // Garante invariante: um Password válido sempre cumpre os critérios
}
```

**Benefícios:**
- Encapsulamento: A validação nunca é esquecida
- Segurança: Impossível ter um Password inválido no sistema
- Coesão: Regras de negócio em um único lugar

### 3. **Padrão Repository**

Implementado para permitir diferentes estratégias de persistência:
```go
type PasswordRepository interface {
    Store(ctx context.Context, password *Password) error
    // Futuras operações: Find, Delete, etc.
}
```

Atualmente usa in-memory, mas pode evoluir para BD facilmente.

### 4. **Pattern Presenter**

Separa a lógica de apresentação da lógica de negócio:
- Output é sempre formatado consistentemente
- Fácil adicionar novos formatos (XML, Protocol Buffers, etc.)

### 5. **Dependency Injection**

Implementado manualmente no router:
```go
passwordRepository := repository.NewPasswordRepository()
presenter := presenter.NewValidatePasswordPresenter()
useCase := usecase.NewValidatePasswordUseCase(duration, repository, presenter)
controller := controller.NewValidatePasswordController(useCase)
```

**Benefícios:**
- Fácil substituir implementações (para testes com mocks)
- Sem magic annotations
- Controle total sobre o grafo de dependências

### 6. **Erros de Domínio**

Tipos de erro customizados para diferentes situações:
```go
type InvalidField struct {
    Field string
    AsIs  string
}
```

**Benefícios:**
- Tratamento de erro específico por camada
- Type assertions claras
- Mensagens de erro em português (ou qualquer idioma)

### 7. **Context e Timeouts**

Propagação de contexto em toda a aplicação:
- Respeita deadlines do cliente
- Facilita cancelamento de operações
- Integração com tracing distribuído

---

## ✅ Testes

A solução possui testes abrangentes em múltiplos níveis:

### Testes de Domínio (Unit Tests)
```bash
go test ./core/domain/password -v
```

Testa o agregado Password com casos como:
- Validações individuais (dígito, maiúscula, etc.)
- Caracteres repetidos
- Espaços em branco
- Senhas válidas e inválidas

### Testes de Caso de Uso (Unit Tests)
```bash
go test ./core/usecase -v
```

Testa a orquestração da lógica usando mocks do repositório.

### Testes de Controller (Integration Tests)
```bash
go test ./adapter/controller -v
```

Testa o fluxo HTTP completo com diferentes cenários de erro.

### Executar Todos os Testes
```bash
go test ./...
```

---

## 📚 Dependências Principais

| Dependência | Versão | Propósito |
|---|---|---|
| gin-gonic/gin | v1.10.0 | Framework HTTP |
| spf13/viper | v1.16.0 | Configuração |
| stretchr/testify | v1.10.0 | Testes (assert, mock) |
| swaggo/swag | v1.8.12 | Documentação Swagger |
| **itau-corp/itau-jw1-dep-golibs-gotel** | **v1.0.2** | **Observabilidade (logs, tracing) - PACOTE INTERNO ITAU** |

### ⚠️ Dependência Interna - Acesso Restrito

A dependência `github.com/itau-corp/itau-jw1-dep-golibs-gotel` é um **pacote interno do Itau** e sua instalação requer acesso ao repositório privado Itau.

#### Pré-requisitos para Download

Para executar este projeto, você precisa:

1. **Estar conectado à rede interna do Itau** ou ter acesso VPN configurado
2. **Ter credenciais de autenticação** configuradas para o Artifactory Itau
3. **Variáveis de ambiente configuradas** conforme o arquivo `.env`:

```dotenv
GOINSECURE='*.prod.aws.cloud.ihf'
GOPROXY=https://artifactory.prod.aws.cloud.ihf/artifactory/api/go/go-remotes,https://artifactory.prod.aws.cloud.ihf/artifactory/itau-jw1-go-release
GONOSUMDB=github.com/itau-corp
```

#### Se Você Não Tiver Acesso

Se você está testando este código **fora do ambiente Itau**, poderá:
- Visualizar o código-fonte
- Executar testes sem conectar ao servidor HTTP (testes unitários apenas)
- Substituir a dependência por uma implementação mock de observabilidade

Porém, **a execução completa da aplicação requer acesso ao repositório interno Itau**.

---

## 🔧 Configuração

Variáveis de ambiente disponíveis (definidas em `.env`):

| Variável | Padrão | Descrição |
|---|---|---|
| APP_NAME | password-validator | Nome da aplicação |
| APP_VERSION | 0.0.1 | Versão da aplicação |
| ENVIRONMENT | local | Ambiente (local, dev, prod) |
| LOGGING_LEVEL | INFO | Nível de log (DEBUG, INFO, WARN, ERROR) |
| HTTP_SERVER_PORT | 8080 | Porta do servidor HTTP |
| SERVER_TIMEOUT | 10 | Timeout em segundos para requisições |
| OTEL_EXPORTER_OTLP_ENDPOINT | http://localhost:4317 | Endpoint do collector OpenTelemetry |

---

## 📊 Observabilidade

A aplicação integra-se com **OpenTelemetry** (via itau-jw1-dep-golibs-gotel) para:

- **Logs estruturados**: Contextualizados com trace IDs
- **Distributed Tracing**: Rastreamento de requisições através das camadas
- **Métricas**: Monitoramento de desempenho (opcional com OpenTelemetry collector)

---

## 🎯 Princípios SOLID Aplicados

### Single Responsibility Principle
- Cada classe tem uma única razão para mudar
- `Password` valida senhas, `Controller` lida com HTTP, etc.

### Open/Closed Principle
- Aberto para extensão: novos tipos de validação podem ser adicionados
- Fechado para modificação: não precisa alterar código existente

### Liskov Substitution Principle
- Implementações de `PasswordRepository` são intercambiáveis
- Implementações de `ValidatePasswordPresenter` são intercambiáveis

### Interface Segregation Principle
- Interfaces pequenas e específicas
- `PasswordRepository` contém apenas operações relevantes

### Dependency Inversion Principle
- Dependências em abstrações (interfaces), não em implementações concretas
- `ValidatePasswordUseCase` depende de `PasswordRepository`, não de uma implementação específica

---

## 🔍 Clean Code

Aplicadas práticas de clean code:

- ✅ **Nomes descritivos**: `ValidatePasswordController`, `InvalidField`
- ✅ **Funções pequenas**: Cada função faz uma coisa bem
- ✅ **Sem código duplicado**: Lógica centralizada no domínio
- ✅ **Error handling explícito**: Erros retornados como valores
- ✅ **Comentários significativos**: Código é auto-documentado quando possível
- ✅ **Formatação consistente**: Go fmt aplicado

---

## 📝 Premises Assumidas

1. **Espaços em branco são ignorados**
   - A especificação menciona que espaços não devem ser considerados como caracteres válidos
   - Implementado removendo espaços antes da validação

2. **Minúsculo vs Maiúsculo com Diacríticos**
   - Unicode é totalmente suportado via `unicode.Is*` do Go
   - Funciona corretamente com caracteres acentuados

3. **Persistência**
   - Atualmente usa in-memory (interface `PasswordRepository`)
   - Facilmente extensível para banco de dados

4. **Concorrência**
   - A aplicação é thread-safe em estado estacionário
   - Contexto propagado para controlar timeouts

5. **Tratamento de Erros**
   - Erros são tipados e tratados de forma granular
   - HTTP 400 para entrada inválida, 500 para erros internos

---

## 🧪 Exemplo de Uso Completo

### 1. Inicie a aplicação:
```bash
go run main.go
```

### 2. Teste um endpoint (usando curl):

```bash
# Senha válida
curl -X POST http://localhost:8080/password/validate \
  -H "Content-Type: application/json" \
  -d '{"password": "AbTp9!fok"}'

# Resposta:
# {"isValid": true}

# Senha inválida
curl -X POST http://localhost:8080/password/validate \
  -H "Content-Type: application/json" \
  -d '{"password": "abc"}'

# Resposta:
# {"errors": [{"field": "password", "message": "Must have at least 9 characters (excluding spaces)"}]}
```

### 3. Visualize a documentação Swagger:
```
http://localhost:8080/swagger/index.html
```

---

## 📦 Estrutura de Resposta

### Sucesso (200 OK)
```json
{
  "isValid": true
}
```

### Erro (400 Bad Request)
```json
{
  "errors": [
    {
      "field": "password",
      "message": "Must contain at least one digit"
    }
  ]
}
```

---

## 🚨 Tratamento de Erros

A aplicação diferencia erros em diferentes camadas:

| Camada | Tipo | Tratamento |
|---|---|---|
| Domínio | `InvalidField` | Retornado ao controller via usecase |
| Application | Parsing | HTTP 400 com detalhes |
| Infrastructure | Servidor | HTTP 500 com log interno |

---

## 💡 Próximas Melhorias Possíveis

- [ ] Adicionar persistência em banco de dados
- [ ] Implementar rate limiting
- [ ] Adicionar autenticação/autorização
- [ ] Caching de resultados
- [ ] Validação assíncrona para senhas
- [ ] Suporte a múltiplas políticas de senha
- [ ] Métricas Prometheus

---

## 📄 Licença

Este projeto foi desenvolvido como parte de um processo seletivo.

---

## 📞 Contato

| Nome             | Email                      |
|------------------|----------------------------|
| Rafael Kassawara | rafael.kassawara@gmail.com |
