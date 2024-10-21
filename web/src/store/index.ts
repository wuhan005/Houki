import { createPinia } from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';

import useAppStore from './modules/app';

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

export { useAppStore };
export default pinia;
