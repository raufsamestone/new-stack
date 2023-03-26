import { useState } from "react";

type Product = {
  Id: number;
  Name: string;
  Price: number;
};

const BASE_URL = "http://localhost:8080";

async function fetchProducts() {
  const response = await fetch(`${BASE_URL}/products`);
  const products = await response.json();
  return products as Product[];
}

async function fetchProductById(id: number) {
  const response = await fetch(`${BASE_URL}/products/${id}`);
  const product = await response.json();
  return product as Product;
}

async function createProduct(product: Product) {
  const response = await fetch(`${BASE_URL}/products`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(product),
  });
  const newProduct = await response.json();
  return newProduct as Product;
}

async function updateProduct(product: Product) {
  const response = await fetch(`${BASE_URL}/products/${product.Id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(product),
  });
  return response.ok;
}

async function deleteProduct(Id: number) {
  const response = await fetch(`${BASE_URL}/products/${Id}`, {
    method: "DELETE",
  });
  return response.ok;
}

export default function Home() {
  const [products, setProducts] = useState<Product[]>([]);

  async function loadProducts() {
    const products = await fetchProducts();
    setProducts(products);
  }

  async function handleCreateProduct() {
    const newProduct = await createProduct({
      name: "New Product",
      price: 0,
    });
    setProducts([...products, newProduct]);
  }

  async function handleUpdateProduct(updatedProduct: Product) {
    await updateProduct(updatedProduct);
    setProducts((prevProducts) =>
      prevProducts.map((product) =>
        product.Id === updatedProduct.Id ? updatedProduct : product
      )
    );
  }

  async function handleDeleteProduct(Id: number) {
    await deleteProduct(Id);
    setProducts((prevProducts) =>
      prevProducts.filter((product) => product.Id !== Id)
    );
  }

  console.log(products);

  return (
    <div className="container mx-auto py-4">
      <div className="flex flex-col space-y-4">
        <button
          className="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded"
          onClick={loadProducts}
        >
          Load Products
        </button>
        <button
          className="bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded"
          onClick={handleCreateProduct}
        >
          Create Product
        </button>
        {products.map((product) => (
          <div key={product.Id} className="flex space-x-4">
            <input
              className="border border-gray-400 rounded px-4 py-2"
              value={product.Name}
              onChange={(e) =>
                handleUpdateProduct({
                  ...product,
                  Name: e.target.value,
                })
              }
            />
            <input
              className="border border-gray-400 rounded px-4 py-2"
              value={product.Price}
              onChange={(e) =>
                handleUpdateProduct({
                  ...product,
                  Price: parseInt(e.target.value),
                })
              }
            />
            <button
              className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded"
              onClick={() => handleDeleteProduct(product.Id)}
            >
              Delete
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
