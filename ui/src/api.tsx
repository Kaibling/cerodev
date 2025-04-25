import { storage_token, storage_api_url, DEFAULT_API_BASE_URL } from "./config";

interface ApiOptions extends Omit<RequestInit, "method" | "body"> {
    method?: "GET" | "POST" | "PUT" | "DELETE";
    body?: Record<string, any>;
}

export async function apiRequest(
    path: string,
    { method = "GET", body, ...rest }: ApiOptions = {}
): Promise<any> {
    const token = localStorage.getItem(storage_token);

    const headers: HeadersInit = {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...rest.headers,
    };

    const options: RequestInit = {
        method,
        headers,
        ...rest,
    };

    if (body) {
        options.body = JSON.stringify(body);
    }
    const CD_API_BASE_URL = localStorage.getItem(storage_api_url);
    if (!CD_API_BASE_URL) {
        localStorage.setItem(storage_api_url, DEFAULT_API_BASE_URL);
    }

    const endpoint = `${CD_API_BASE_URL}${path}`
    const response = await fetch(endpoint, options);

    if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(error.message || "API Error");
    }

    return response.json();
}



