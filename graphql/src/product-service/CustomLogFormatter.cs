using System.Text.Encodings.Web;
using System.Text.Json;
using Serilog.Events;
using Serilog.Formatting;

// Probably faster to direct use output.Write like in
// https://github.com/serilog/serilog/blob/dev/src/Serilog/Formatting/Json/JsonFormatter.cs#L55
// but this is for simplicity and readability
// perhaps we benchmark later?
public class CustomLogFormatter : ITextFormatter
{
    public CustomLogFormatter() { }

    private static readonly JsonSerializerOptions _jsonOptions = new JsonSerializerOptions
    {
        Encoder = JavaScriptEncoder.UnsafeRelaxedJsonEscaping,
    };

    private static readonly string[] _allowedProperties =
    {
        "RequestPath", "SourceContext"
    };

    public void Format(LogEvent logEvent, TextWriter output)
    {
        var log = new Dictionary<string, string>()
        {
            {"level", logEvent.Level.ToString().ToLower()},
            {"timestamp", logEvent.Timestamp.UtcDateTime.ToString("O")},
            {"message", logEvent.MessageTemplate.Text}
        };

        if (logEvent.TraceId != null)
        {
            log.Add("trace.id", logEvent.TraceId.ToString()!);
        }

        if (logEvent.SpanId != null)
        {
            log.Add("span.id", logEvent.SpanId.ToString()!);
        }

        logEvent.RemovePropertyIfPresent("ConnectionId");

        foreach (var property in logEvent.Properties.Where(x => _allowedProperties.Contains(x.Key, StringComparer.OrdinalIgnoreCase)))
        {
            log.Add(property.Key, property.Value.ToString());
        }

        string logLine = JsonSerializer.Serialize(log, _jsonOptions);

        output.Write(logLine);
        output.Write(Environment.NewLine);
    }
}