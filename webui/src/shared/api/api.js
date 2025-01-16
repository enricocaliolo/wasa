import axios from 'axios';
import { useUserStore } from '../stores/user_store';

const api = axios.create({
	baseURL: __API_URL__,
	timeout: 5000,
	headers: {
		'Content-Type': 'application/json',
		'Access-Control-Allow-Origin': '*',
	},
});

api.interceptors.request.use((config) => {
	const userStore = useUserStore();

	if (userStore.user.userId !== undefined) {
		config.headers.Authorization = `Bearer ${userStore.user.userId}`;
	}
	return config;
});

export default api;
