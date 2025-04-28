import { Play, Trash2, Square, ExternalLink } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Plus } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { type ReactElement } from 'react';
import { apiRequest } from '@/api';
import { User } from '../App';
import { get_base_api_url } from '@/config';

interface Container {
    id: string;
    docker_id: string;
    image_name: string;
    status: string;
    state: string;
    container_name: string;
    git_repo: string;
    user_id: string;
    env_vars: string[];
    ports: string[];
    ui_port: string;
}
interface ContainerProps {
    user: User | null;
}


const statusColor: Record<string, string> = {
  running: 'text-green-400',
  stopped: 'text-red-400',
};

export default function Containers({ user }: ContainerProps): ReactElement {
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [containerList, setContainerList] = useState<Container[]>([]);
  const [newContainer, setNewContainer] = useState({
    id: '',
    docker_id: '',
    image_name: '',
    status: '',
    state: '',
    container_name: '',
    git_repo: '',
    user_id: '',
    env_vars: [''],
    ports: [''],
  });
  const [showForm, setShowForm] = useState(false);
  const [availableImages, setAvailableImages] = useState<string[]>([]);

  useEffect(() => {
    apiRequest('/api/v1/images')
      .then((data) => {
        if (data?.data?.length > 0) {
          const formatted = data.data.map((img: any) => `${img.repo_name}:${img.tag}`);
          setAvailableImages(formatted);
        }
      })
      .catch((err) => {
        console.error('Error fetching images:', err);
      });
  }, []);

  const getContainers = () => {
    apiRequest('/api/v1/containers')
      .then((data) => { setContainerList(data.data); })
      .catch((err) => {
        console.error('Error fetching containers:', err);
        setErrorMsg('Could not load containers. Please try again later.');
      });
  };


  useEffect(() => {
    getContainers();
  }, []);

  const handleAdd = async () => {
    if (!newContainer.image_name || !newContainer.git_repo || !newContainer.ports) return;

    apiRequest('/api/v1/containers', {
      method: 'POST',
      body: { ...newContainer, user_id: user?.user_id },
    }).then((data) => {
      setNewContainer({
        id: '',
        docker_id: '',
        image_name: '',
        status: '',
        state: '',
        container_name: '',
        git_repo: '',
        user_id: '',
        env_vars: [''],
        ports: [''],
      });
      setShowForm(false);
      getContainers();
    })
      .catch((err) => {
        console.error('Error fetching containers:', err);
        setErrorMsg('Could not load containers. Please try again later.');
      });
  };

  const handleDelete = async (container_id: string) => {
    const confirmDelete = window.confirm('Are you sure you want to delete this container?');
    if (!confirmDelete) return;
    apiRequest(`/api/v1/containers/${container_id}`, { method: 'DELETE' }).then((data) => {
      getContainers();
    }).catch((err) => {
      console.error('Error deleting containers:', err);
      setErrorMsg('Could not load containers. Please try again later.');
    });
    getContainers();
  };

  const handleStart = async (container_id: string) => {
    apiRequest(`/api/v1/containers/${container_id}/start`, { method: 'POST' }).then((data) => {
      getContainers();
    }).catch((err) => {
      console.error('Error starting containers:', err);
      setErrorMsg('Could not load containers. Please try again later.');
    });
  };

  return (
    <main className="p-6">
      {errorMsg && (
        <div className="mb-4 p-3 bg-red-600 text-white rounded-xl">
          {errorMsg}
        </div>
      )}
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-xl font-semibold">Your Containers</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="p-2 bg-green-600 hover:bg-green-500 rounded-xl"
        >
          <Plus className="w-5 h-5 text-white" />
        </button>
      </div>

      <AnimatePresence>
        {showForm && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="overflow-hidden mb-8"
          >
            <div className="bg-gray-800 p-4 rounded-xl border border-gray-700">
              <h3 className="text-white font-semibold mb-2">Add New Container</h3>
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <select
                  className="p-2 rounded bg-gray-700 text-white"
                  value={newContainer.image_name}
                  onChange={(e) => setNewContainer({ ...newContainer, image_name: e.target.value })}
                >
                  <option value="">Select Image</option>
                  {availableImages.map((img) => (
                    <option key={img} value={img}>{img}</option>
                  ))}
                </select>
                <input
                  className="p-2 rounded bg-gray-700 text-white"
                  placeholder="Repo"
                  value={newContainer.git_repo}
                  onChange={(e) => setNewContainer({ ...newContainer, git_repo: e.target.value })}
                />
                <div className="col-span-full">
                  <label className="text-sm font-medium text-white mb-1 block">Ports</label>
                  {newContainer.ports.map((port, idx) => (
                    <div key={idx} className="flex gap-2 mb-2">
                      <input
                        className="flex-1 p-2 rounded bg-gray-700 text-white"
                        placeholder="8097:8080/tcp"
                        value={port}
                        onChange={(e) => {
                          const newPorts = [...newContainer.ports];
                          newPorts[idx] = e.target.value;
                          setNewContainer({ ...newContainer, ports: newPorts });
                        }}
                      />
                      <button
                        type="button"
                        onClick={() => {
                          const newPorts = newContainer.ports.filter((_, i) => i !== idx);
                          setNewContainer({ ...newContainer, ports: newPorts });
                        }}
                        className="bg-red-600 hover:bg-red-500 text-white px-2 rounded"
                      >
                                                −
                      </button>
                    </div>
                  ))}
                  <button
                    type="button"
                    onClick={() =>
                      setNewContainer({ ...newContainer, ports: [...newContainer.ports, ''] })
                    }
                    className="mt-2 px-3 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded"
                  >
                                        + Add Port
                  </button>
                </div>
                <div className="col-span-full">
                  <label className="text-sm font-medium text-white mb-1 block">Environment Variables</label>
                  {newContainer.env_vars.map((env, idx) => (
                    <div key={idx} className="flex gap-2 mb-2">
                      <input
                        className="flex-1 p-2 rounded bg-gray-700 text-white"
                        placeholder="ENV_VAR=value"
                        value={env}
                        onChange={(e) => {
                          const newEnvs = [...newContainer.env_vars];
                          newEnvs[idx] = e.target.value;
                          setNewContainer({ ...newContainer, env_vars: newEnvs });
                        }}
                      />
                      <button
                        type="button"
                        onClick={() => {
                          const newEnvs = newContainer.env_vars.filter((_, i) => i !== idx);
                          setNewContainer({ ...newContainer, env_vars: newEnvs });
                        }}
                        className="bg-red-600 hover:bg-red-500 text-white px-2 rounded"
                      >
                                                -
                      </button>
                    </div>
                  ))}
                  <button
                    type="button"
                    onClick={() =>
                      setNewContainer({ ...newContainer, env_vars: [...newContainer.env_vars, ''] })
                    }
                    className="mt-2 px-3 py-1 bg-blue-600 hover:bg-blue-500 text-white rounded"
                  >
                                        + Add Env Var
                  </button>
                </div>
              </div>
              <button
                className="mt-4 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-xl flex items-center gap-2"
                onClick={handleAdd}
              >
                <Plus className="w-4 h-4" /> Add Container
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="space-y-4">
        {containerList.length === 0 ? (
          <div className="text-center text-gray-400 italic">No containers running.</div>
        ) : (
          containerList.map((c: Container, index: number) => (
            <div
              key={index}
              className="flex justify-between items-start bg-gray-800 rounded-xl p-4 border border-gray-700"
            >
              <div>
                <h3>{c.container_name.split('-').slice(2).join('-')}</h3>
                <div className="text-sm space-y-1 mt-1">
                  <div className={`font-medium ${statusColor[c.state]}`}>{c.state}</div>
                  <div className="text-gray-400">{c.status}</div>
                  <div className="text-gray-400">{c.image_name}</div>
                  <div className="text-gray-400">{c.git_repo}</div>
                  <div className="text-gray-400">Port: {c.ports}</div>
                </div>
              </div>
              <div className="flex gap-2">
                {c.state === 'running' ? (
                  <button className="p-2 bg-blue-600 hover:bg-blue-500 rounded-xl">
                    <Square className="w-5 h-5" />
                  </button>
                ) : (
                  <button className="p-2 bg-blue-600 hover:bg-blue-500 rounded-xl" onClick={() => handleStart(c.id)}>
                    <Play className="w-5 h-5" />
                  </button>
                )}
                <button className="p-2 bg-red-600 hover:bg-red-500 rounded-xl" onClick={() => handleDelete(c.id)}>
                  <Trash2 className="w-5 h-5" />
                </button>
                {c.state === 'running' && (
                  <a
                    href={`${get_base_api_url()}/proxy/${c.id}`}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    <button className="p-2 bg-gray-700 hover:bg-gray-600 rounded-xl">
                      <ExternalLink className="w-5 h-5" />
                    </button>
                  </a>
                )}
              </div>
            </div>
          )))}
      </div>
    </main>
  );
}