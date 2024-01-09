using System.Text.Json;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;

public class OrdersConsumer : BackgroundService
{
    private readonly ILogger<OrdersConsumer> _logger;
    private readonly IServiceProvider _serviceProvider;
    private string? rabbitMQServer = Environment.GetEnvironmentVariable("RABBITMQ_HOST");
    private string? rabbitMQPort = Environment.GetEnvironmentVariable("RABBITMQ_PORT");
    private string? rabbitMQUser = Environment.GetEnvironmentVariable("RABBITMQ_USERNAME");
    private string? rabbitMQPass = Environment.GetEnvironmentVariable("RABBITMQ_PASSWORD");
    private readonly string queueName = "orders";

    public OrdersConsumer(ILogger<OrdersConsumer> logger, IServiceProvider serviceProvider)
    {
        _logger = logger;
        _serviceProvider = serviceProvider;
    }

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        stoppingToken.ThrowIfCancellationRequested();

        var db = _serviceProvider.CreateScope().ServiceProvider.GetRequiredService<OrderContext>();

        int port;
        int.TryParse(rabbitMQPort, out port);

        var factory = new ConnectionFactory()
        {
            HostName = rabbitMQServer,
            Port = port,
            UserName = rabbitMQUser,
            Password = rabbitMQPass
        };

        try
        {
            var connection = factory.CreateConnection();
            var channel = connection.CreateModel();

            var c = channel.QueueDeclare(queue: queueName,
                                    durable: true,
                                    exclusive: false,
                                    autoDelete: false,
                                    arguments: null);

            var consumer = new EventingBasicConsumer(channel);

            _logger.LogInformation("READY TO RECEIVE");

            consumer.Received += async (model, ea) =>
            {
                _logger.LogInformation("RECEIVED");

                var body = ea.Body.ToArray();
                try
                {
                    var request = JsonSerializer.Deserialize<Order>(body);
                    if (request is null)
                    {
                        return;
                    }

                    

                    await db.Orders.AddAsync(request);
                    await db.SaveChangesAsync();
                }
                catch (Exception ex)
                {
                    _logger.LogError($"Error while processing order: {ex.Message} {ex.StackTrace}");
                }
            };

            channel.BasicConsume(queue: c.QueueName,
                autoAck: true,
                consumer: consumer);

            await Task.CompletedTask;
        }
        catch (Exception ex)
        {
            _logger.LogError($"Error while connecting to RabbitMQ: '{ex.Message}' {ex.StackTrace}");
        }
    }
}