import AdminCreateProductForm from "./components/create-product-form";

export default function AdminCreateProductPage() {
  return (
    <>
      <h1 className="text-xl mb-2">Create your product</h1>

      <div className="w-1/2">
        <AdminCreateProductForm />
      </div>
    </>
  );
}
