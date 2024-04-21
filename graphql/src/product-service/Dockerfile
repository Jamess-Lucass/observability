FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS builder
WORKDIR /app

# Copy everything
COPY . ./
# Restore as distinct layers
RUN dotnet restore
# Build and publish a release
RUN dotnet publish -c Release -o out

# Build runtime image
FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine
WORKDIR /app
COPY --from=builder /app/out .

ENTRYPOINT ["dotnet", "product-service.dll"]