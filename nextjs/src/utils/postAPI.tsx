export async function postAPI<T>(url: string, data: any, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    ...options,
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(options?.headers || {}),
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    throw new Error(`An error occurred: ${response.status}`);
  }
  const responseData = await response.json();
  return responseData as T;
}
