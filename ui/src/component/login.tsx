import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { storage_token } from "../config";
import { type ReactElement } from "react";
import { type User } from "../App";
import { apiRequest } from "@/api";

interface LoginPageProps {
    setUser: React.Dispatch<React.SetStateAction<User | null>>;

}

export default function LoginPage({ setUser }: LoginPageProps): ReactElement {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate();

    const handleLogin = async () => {
        apiRequest(`/api/v1/auth/login`, { method: "POST", headers: { "Content-Type": "application/json" }, body: { username, password } })
            .then((data) => {
                localStorage.setItem(storage_token, data.data.token);
                setUser({ user_id: data.data.user_id, username: data.data.username });
                navigate("/containers");
            })
            .catch((err) => {
                console.error(err);
                alert("Login failed");
            });
    };

    const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            handleLogin();
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white">
            <div className="bg-gray-800 p-8 rounded-xl w-full max-w-sm">
                <h2 className="text-xl font-bold mb-4">Login</h2>
                <input
                    className="w-full mb-3 p-2 rounded bg-gray-700"
                    type="username"
                    placeholder="Username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                />
                <input
                    className="w-full mb-4 p-2 rounded bg-gray-700"
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    onKeyDown={handleKeyPress}
                />
                <button
                    onClick={handleLogin}
                    className="w-full bg-green-600 hover:bg-green-500 p-2 rounded"
                >
                    Sign In
                </button>
            </div>
        </div>
    );
}
