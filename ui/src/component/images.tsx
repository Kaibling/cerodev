import { useEffect, useState } from "react";
import { type ReactElement } from "react";
import { apiRequest } from "@/api";

interface Image {
    image_id: string;
    repo_name: string;
    tag: string;
}

export default function Images(): ReactElement {
    const [errorMsg, setErrorMsg] = useState<string | null>(null);
    const [imageList, setImageList] = useState<Image[]>([]);


    useEffect(() => {
        apiRequest(`/api/v1/images`)
            .then((data) => {
                setImageList(data.data);
            })
            .catch((err) => {
                console.error("Error fetching images:", err);
                setErrorMsg("Could not load images. Please try again later.");
            });
    }, []);


    return (
        <main className="p-6">
            {errorMsg && (
                <div className="mb-4 p-3 bg-red-600 text-white rounded-xl">
                    {errorMsg}
                </div>
            )}
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-xl font-semibold">Your Images</h2>
            </div>

            <div className="space-y-4">
                {imageList.length === 0 ? (
                    <div className="text-center text-gray-400 italic">No Images build.</div>
                ) : (
                    imageList.map((i: Image, index: number) => (
                        <div
                            key={index}
                            className="flex justify-between items-start bg-gray-800 rounded-xl p-4 border border-gray-700"
                        >
                            <div>
                                <h3 className="text-lg font-semibold text-white">{i.repo_name.split('-').pop()}</h3>
                                <div className="text-sm space-y-1 mt-1">
                                    <div className="text-gray-400">tag: {i.tag}</div>
                                    <div className="text-gray-400">id: {i.image_id}</div>
                                </div>
                            </div>
                        </div>
                    )))}
            </div>
        </main>
    );
}