using FluentValidation;

public class OrderValidator : AbstractValidator<Order>
{
    public OrderValidator(IProductsClient productsClient)
    {
        RuleFor(x => x.ProductId).NotEmpty().MustAsync(async (id, cancellation) =>
        {
            var response = await productsClient.GetAsync(id);

            if (response is null)
            {
                return false;
            }

            return true;
        }).WithMessage("Invalid ProductId");
        RuleFor(x => x.Quantity).NotEmpty().GreaterThan(0u).LessThan(1_000_000u);
        RuleFor(x => x.Price).NotEmpty().PrecisionScale(18, 2, false);
    }
}