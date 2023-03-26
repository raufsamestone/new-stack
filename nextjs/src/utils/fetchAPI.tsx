export async function fetchAPI<T>(
  url: string,
  options?: RequestInit
): Promise<T> {
  const response = await fetch(url, options);
  if (!response.ok) {
    throw new Error(`An error occurred: ${response.status}`);
  }
  const data = await response.json();
  return data as T;
}
