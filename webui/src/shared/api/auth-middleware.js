import { User } from '../../modules/auth/models/user';
import { useUserStore } from '../stores/user_store';
import api from './api';

export async function authMiddleware(to, from, next) {
	const userStore = useUserStore();

	if (to.name === 'login') {
		next();
		return;
	}

	const savedUsername = localStorage.getItem('username');
	if (savedUsername) {
		if (userStore.user.userId === undefined) {
			try {
				const response = await api.put('/session', {
					username: savedUsername,
				});

				if (response.data) {
					userStore.setUser(User.fromJSON(response.data));
					next();
					return;
				}
			} catch (error) {
				console.error('Auto-login failed:', error);
				localStorage.removeItem('username');
				next('/login');
				return;
			}
		} else {
			next();
			return;
		}
	}

	if (userStore.user.userId === undefined) {
		next('/login');
	} else {
		next();
	}
}
