# 👨‍🍳 GoLunch Operation Service

Microsserviço responsável pelo gerenciamento das operações da cozinha e painel administrativo da lanchonete GoLunch. Este serviço gerencia a fila de pedidos para a cozinha, atualizações de status e controle administrativo.

## 🎯 Responsabilidades

- **Gestão de Operações**: Controle da fila de pedidos para a cozinha
- **Atualização de Status**: Gerenciamento do fluxo de status dos pedidos
- **Painel Administrativo**: Interface para administradores gerenciarem o sistema
- **Autenticação**: Sistema de login e autorização para administradores
- **Monitoramento**: Acompanhamento de tempo de preparo e status

## 🏗️ Arquitetura

O serviço segue os princípios da **Arquitetura Hexagonal** com as seguintes camadas:

- **Entities**: Regras de negócio fundamentais
- **Use Cases**: Lógica de negócio específica
- **Gateways**: Interfaces para acesso a dados externos
- **Controllers**: Coordenação entre camadas
- **Handlers**: Gerenciamento de requisições HTTP
- **External/Infrastructure**: Implementações concretas (banco de dados)

## 🗄️ Banco de Dados

- **PostgreSQL**: Banco de dados principal
- **Tabelas**:
  - `admins`: Dados dos administradores
  - `orders`: Pedidos (read-only, sincronizado com Core Service)

## 🚀 Endpoints Disponíveis

### Autenticação
- `POST /admin/register` - Cadastrar novo administrador
- `POST /admin/login` - Login de administrador

### Gestão de Pedidos (Admin)
- `GET /admin/orders` - Listar todos os pedidos
- `PUT /admin/orders/:id` - Atualizar status do pedido
- `GET /admin/orders/panel` - Painel de pedidos para cozinha

### Health Check
- `GET /ping` - Health check do serviço

## 🔧 Configuração Local

1. **Clone o repositório**
2. **Configure as variáveis de ambiente**:
   ```bash
   export DATABASE_URL="postgres://user:password@localhost:5432/golunch_operation?sslmode=disable"
   export SECRET_KEY="your-jwt-secret-key"
   ```

3. **Execute o banco de dados**:
   ```bash
   docker-compose up -d postgres
   ```

4. **Execute a aplicação**:
   ```bash
   go run cmd/api/main.go
   ```

## 📋 Dependências

- **Go** 1.24.3
- **PostgreSQL** 16.3
- **Gin** - Framework web
- **GORM** - ORM para banco de dados
- **JWT** - Autenticação e autorização
- **Swagger** - Documentação da API

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes BDD
go test -tags=bdd ./...
```

## 📊 Cobertura de Testes

- **Meta**: 80% de cobertura
- **BDD**: Implementado para cenários de atualização de status
- **Testes Unitários**: Todos os use cases e controllers
- **Testes de Integração**: Autenticação e autorização

## 🐳 Docker

```bash
# Build da imagem
docker build -t tc-golunch-operation-service .

# Executar container
docker run -p 8083:8083 tc-golunch-operation-service
```

## 📈 Monitoramento

- **Health Check**: `GET /ping`
- **Swagger UI**: `GET /swagger/index.html`
- **Logs**: Estruturados em JSON
- **Métricas**: Tempo de preparo, status de pedidos

## 🔄 CI/CD

O serviço possui pipeline CI/CD configurado com:
- Validação de código
- Execução de testes
- Análise de cobertura
- Build e deploy automático
- Proteção de branch main

## 🔐 Segurança

- **JWT Tokens**: Autenticação segura
- **Middleware de Autorização**: Controle de acesso
- **Admin Only**: Endpoints restritos a administradores
- **HTTPS**: Comunicação segura

## 📝 Documentação da API

A documentação completa da API está disponível via Swagger UI em:
`http://localhost:8083/swagger/index.html`

## 🔗 Integração com Outros Serviços

- **Core Service**: Sincronização de pedidos
- **Payment Service**: Notificações de status de pagamento
- **Message Queue**: Comunicação assíncrona entre serviços

## 👥 Fluxo de Trabalho da Cozinha

1. **Pedido Recebido**: Pedido aparece na fila da cozinha
2. **Em Preparação**: Administrador marca como "em preparação"
3. **Pronto**: Administrador marca como "pronto"
4. **Finalizado**: Pedido é marcado como "finalizado" após retirada

## 📊 Painel Administrativo

- **Fila de Pedidos**: Lista de pedidos pendentes
- **Status em Tempo Real**: Atualizações instantâneas
- **Tempo de Preparo**: Controle de tempo estimado
- **Histórico**: Relatórios de operações da cozinha

## 🔗 Integração Serverless (AWS Lambda)

✅ A autenticação serverless já está configurada.

### **🛠️ Código Implementado**
O código foi atualizado seguindo o padrão do monolítico `tc-golunch-api`:

1. **ServerlessAuthGateway**: Implementado para comunicação com Lambda
2. **ServerlessAuthMiddleware**: Middleware de autenticação serverless
3. **ServerlessAdminOnly**: Middleware específico para validação de admin via serverless
4. **main.go**: Atualizado para usar serverless auth em vez de JWT local

### **🔧 Configuração das URLs**

**⚠️ PREREQUISITO**: Primeiro faça deploy do `tc-golunch-serverless` para gerar as URLs reais!

```bash
# 1. Deploy serverless (OBRIGATÓRIO primeiro)
cd ../tc-golunch-serverless
terraform init
terraform apply
# Isso cria funções Lambda e gera URLs reais do API Gateway

# 2. Obter URLs reais geradas
terraform output
# Output: api_gateway_url = "https://abc123def.execute-api.us-east-1.amazonaws.com"

# 3. ENTÃO configurar variáveis locais com URLs reais:
export LAMBDA_AUTH_URL="https://abc123def.execute-api.us-east-1.amazonaws.com/auth"
export SERVICE_AUTH_LAMBDA_URL="https://abc123def.execute-api.us-east-1.amazonaws.com/service-auth"

# Variáveis existentes (mantidas)
export DATABASE_URL="host=localhost user=golunch_prod password=golunch_prod123 dbname=golunch_production port=5434 sslmode=disable TimeZone=America/Sao_Paulo"
export SECRET_KEY="production-secret-key-2024"
export OPERATION_SERVICE_PORT="8083"
export ORDER_SERVICE_URL="http://localhost:8081"
export PAYMENT_SERVICE_URL="http://localhost:8082"
```

### **📦 Deploy Kubernetes**

⚠️ **PREREQUISITO**: Deploy do `tc-golunch-serverless` ANTES de fazer deploy Kubernetes!

**Passo-a-passo completo:**

```bash
# PASSO 1: Deploy Serverless (OBRIGATÓRIO primeiro)
cd ../tc-golunch-serverless
terraform init
terraform apply

# PASSO 2: Obter URLs reais do API Gateway
terraform output
# Exemplo output: api_gateway_url = "https://abc123def.execute-api.us-east-1.amazonaws.com"

# PASSO 3: Atualizar ConfigMap com URLs REAIS
cd ../tc-golunch-operation-service
vim k8s/operation-service-configmap.yaml

# SUBSTITUIR estas linhas (são templates):
# LAMBDA_AUTH_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/auth"
# SERVICE_AUTH_LAMBDA_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/service-auth"

# POR URLs reais obtidas no terraform output:
# LAMBDA_AUTH_URL: "https://abc123def.execute-api.us-east-1.amazonaws.com/auth"  
# SERVICE_AUTH_LAMBDA_URL: "https://abc123def.execute-api.us-east-1.amazonaws.com/service-auth"

# PASSO 4: Deploy Kubernetes
kubectl apply -f k8s/
```

**Estrutura já configurada:**
```yaml
# k8s/operation-service-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: operation-service-config
data:
  LAMBDA_AUTH_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/auth"
  SERVICE_AUTH_LAMBDA_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/service-auth"
  # ... outras variáveis
```

### **✅ Verificação da Configuração**

Após configurar as variáveis, teste a integração:

```bash
# 1. Inicie o serviço
go run cmd/api/main.go

# 2. Teste login de admin via serverless
curl -X POST http://localhost:8083/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 3. Teste endpoint protegido (requer admin)
curl -X GET http://localhost:8083/admin/orders \
  -H "Authorization: Bearer <token-do-lambda>"

# 4. Verifique logs para confirmação da integração Lambda
```

### **🔄 Migração Gradual**

A implementação mantém **compatibilidade total** com o código existente:
- ✅ Mesmas interfaces de autenticação  
- ✅ Mesmos endpoints e responses
- ✅ Zero breaking changes para clientes
- ✅ Fallback automático se Lambda não disponível
- ✅ **ServerlessAdminOnly** específico para operações administrativas


