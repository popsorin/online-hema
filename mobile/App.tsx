/**
 * HEMA Lessons Mobile App
 *
 * Root component that sets up providers and navigation.
 */

import React from 'react';
import {SafeAreaProvider} from 'react-native-safe-area-context';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import {StatusBar} from 'react-native';
import {RootNavigator} from '@/navigation';

// Create a client for React Query
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      retry: 2,
      refetchOnWindowFocus: false,
    },
  },
});

function App(): React.JSX.Element {
  return (
    <QueryClientProvider client={queryClient}>
      <SafeAreaProvider>
        <StatusBar barStyle="dark-content" backgroundColor="#f5f5f5" />
        <RootNavigator />
      </SafeAreaProvider>
    </QueryClientProvider>
  );
}

export default App;
