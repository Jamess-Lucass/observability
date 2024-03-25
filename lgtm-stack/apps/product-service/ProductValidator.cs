using FluentValidation;

public class ProductValidator : AbstractValidator<Product>
{
    public ProductValidator()
    {
        RuleFor(x => x.Name).Length(2, 128);
        RuleFor(x => x.Description).Length(2, 1024);
        RuleFor(x => x.Price).GreaterThan(0).PrecisionScale(18, 2, false);
    }
}