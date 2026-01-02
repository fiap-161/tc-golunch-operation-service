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


