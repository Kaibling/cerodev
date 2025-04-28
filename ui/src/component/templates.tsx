import { useEffect, useState } from 'react';
import { Trash2, Edit, Plus, Hammer } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import Editor from '@monaco-editor/react';
import { type ReactElement } from 'react';
import { apiRequest } from '@/api';

interface Template {
    id: string;
    name: string;
    repo_name: string;
    dockerfile: string;
}

export default function Templates(): ReactElement {
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editInputs, setEditInputs] = useState<Record<string, Template>>({});
  const [activeBuildForm, setActiveBuildForm] = useState<string | null>(null);
  const [buildInputs, setBuildInputs] = useState<Record<string, { tag: string; repo_name: string }>>({});
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [templates, setTemplates] = useState<Template[]>([]);
  const [newTemplate, setNewTemplate] = useState({
    id: '',
    name: '',
    repo_name: '',
    dockerfile: '',
  });
  const [showForm, setShowForm] = useState(false);

  const getTemplates = async () => {
    apiRequest('/api/v1/templates')
      .then((data) => setTemplates(data.data))
      .catch((err) => {
        console.error('Error fetching templates:', err);
        setErrorMsg('Could not load templates. Please try again later.');
      });
  };

  useEffect(() => {
    getTemplates();
  }, []);

  const handleAdd = async () => {
    if (!newTemplate.name || !newTemplate.repo_name) return;
    apiRequest('/api/v1/templates', {
      method: 'POST',
      body: newTemplate,
    }).then((data) => {
      setNewTemplate({ id: '', name: '', repo_name: '', dockerfile: '' });
      setShowForm(false);
      getTemplates();
    }).catch((err) => {
      console.error('Error fetching templates:', err);
      setErrorMsg('Could not load templates. Please try again later.');
    });

  };

  const handleEditSave = async (id: string) => {
    const updatedTemplate = editInputs[id];
    apiRequest(`/api/v1/templates/${id}`, {
      method: 'PUT',
      body: updatedTemplate,
    }).then((data) => {
      const updated = templates.map((tpl) =>
        tpl.id === id ? updatedTemplate : tpl,
      );
      setTemplates(updated);
      setEditingId(null);
    }).catch((err) => {
      console.error('Error updating template:', err);
      setErrorMsg('Could not update template. Please try again.');
    });
  };

  const handleDelete = async (id: string) => {
    const confirmDelete = window.confirm('Are you sure you want to delete this template?');
    if (!confirmDelete) return;
    apiRequest(`/api/v1/templates/${id}`, {
      method: 'DELETE',
    }).then(() => {
      setTemplates(templates.filter((tpl) => tpl.id !== id));
    }).catch((err) => {
      console.error('Error deleting template:', err);
      setErrorMsg('Could not delete template. Please try again later.');
    });
  };

  const handleBuild = async (id: string, repo_name: string) => {
    const input = buildInputs[id];
    if (!input?.tag) return;

    apiRequest(`/api/v1/templates/${id}`, {
      method: 'POST',
      body: { ...input, repo_name: repo_name },
    }).then((data) => {
      setActiveBuildForm(null);
      setBuildInputs((prev) => {
        const updated = { ...prev };
        delete updated[id];
        return updated;
      });
    }).catch((err) => {
      console.error('Error building image:', err);
      setErrorMsg('Could not start build. Please try again later.');
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
        <h2 className="text-xl font-semibold">Templates</h2>
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
              <h3 className="text-white font-semibold mb-2">Add New Template</h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <input
                  className="p-2 rounded bg-gray-700 text-white"
                  placeholder="Name"
                  value={newTemplate.name}
                  onChange={(e) => setNewTemplate({ ...newTemplate, name: e.target.value })}
                />
                <input
                  className="p-2 rounded bg-gray-700 text-white"
                  placeholder="Repo Name"
                  value={newTemplate.repo_name}
                  onChange={(e) => setNewTemplate({ ...newTemplate, repo_name: e.target.value })}
                />
              </div>
              <button
                className="mt-4 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-xl flex items-center gap-2"
                onClick={handleAdd}
              >
                <Plus className="w-4 h-4" /> Add Template
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="space-y-4">
        {templates.length === 0 ? (
          <div className="text-center text-gray-400 italic">No Templates.</div>
        ) : (
          templates.map((tpl) => (
            <div
              key={tpl.id}
              className="bg-gray-800 p-4 rounded-xl border border-gray-700 space-y-2"
            >
              <div className="flex justify-between items-start">
                <div>
                  <div className="text-lg font-medium text-white">{tpl.name}</div>
                  <div className="text-gray-400 text-sm">Image Name: {tpl.repo_name}</div>
                </div>
                <div className="flex gap-2">
                  <button
                    className="p-2 bg-yellow-600 hover:bg-yellow-500 rounded-xl"
                    onClick={() => {
                      setEditingId(tpl.id === editingId ? null : tpl.id);
                      setEditInputs((prev) => ({
                        ...prev,
                        [tpl.id]: { ...tpl },
                      }));
                    }}
                  >
                    <Edit className="w-5 h-5" />
                  </button>
                  <button
                    className="p-2 bg-red-600 hover:bg-red-500 rounded-xl"
                    onClick={() => handleDelete(tpl.id)}
                  >
                    <Trash2 className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => setActiveBuildForm(activeBuildForm === tpl.id ? null : tpl.id)}
                    className="p-2 bg-green-600 hover:bg-green-500 rounded-xl"
                  >
                    <Hammer className="w-5 h-5" />
                  </button>
                </div>
              </div>

              <AnimatePresence>
                {activeBuildForm === tpl.id && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                    className="overflow-hidden pl-1 pr-1"
                  >
                    <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-2">
                      {/* <input
                                                className="p-2 rounded bg-gray-700 text-white"
                                                placeholder="Image Name"
                                                value={buildInputs[tpl.id]?.repo_name || ""}
                                                onChange={(e) =>
                                                    setBuildInputs((prev) => ({
                                                        ...prev,
                                                        [tpl.id]: { ...prev[tpl.id], repo_name: e.target.value },
                                                    }))
                                                }
                                            /> */}
                      <input
                        className="p-2 rounded bg-gray-700 text-white"
                        placeholder="Tag"
                        value={buildInputs[tpl.id]?.tag || ''}
                        onChange={(e) =>
                          setBuildInputs((prev) => ({
                            ...prev,
                            [tpl.id]: { ...prev[tpl.id], tag: e.target.value },
                          }))
                        }
                      />
                    </div>
                    <button
                      onClick={() => handleBuild(tpl.id, tpl.repo_name)}
                      className="mt-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-xl"
                    >
                                            🚀 Build Image
                    </button>
                  </motion.div>
                )}
              </AnimatePresence>

              <AnimatePresence>
                {editingId === tpl.id && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    exit={{ opacity: 0, height: 0 }}
                    className="overflow-hidden mt-4 space-y-4"
                  >
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                      <input
                        className="p-2 rounded bg-gray-700 text-white"
                        placeholder="Template Name"
                        value={editInputs[tpl.id]?.name || ''}
                        onChange={(e) =>
                          setEditInputs((prev) => ({
                            ...prev,
                            [tpl.id]: { ...prev[tpl.id], name: e.target.value },
                          }))
                        }
                      />
                      <input
                        className="p-2 rounded bg-gray-700 text-white"
                        placeholder="Image Name"
                        value={editInputs[tpl.id]?.repo_name || ''}
                        onChange={(e) =>
                          setEditInputs((prev) => ({
                            ...prev,
                            [tpl.id]: { ...prev[tpl.id], repo_name: e.target.value },
                          }))
                        }
                      />
                    </div>

                    <div className="h-128 border border-gray-700 rounded">
                      <Editor
                        height="100%"
                        defaultLanguage="dockerfile"
                        theme="vs-dark"
                        value={editInputs[tpl.id]?.dockerfile || ''}
                        onChange={(value) =>
                          setEditInputs((prev) => ({
                            ...prev,
                            [tpl.id]: { ...prev[tpl.id], dockerfile: value || '' },
                          }))
                        }
                      />
                    </div>

                    <button
                      onClick={() => handleEditSave(tpl.id)}
                      className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-xl"
                    >
                                            💾 Save Changes
                    </button>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          ))
        )}
      </div>
    </main>
  );
}
