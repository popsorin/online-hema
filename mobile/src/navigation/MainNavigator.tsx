/**
 * Main Navigator
 *
 * Navigation stack for all app screens (content is publicly accessible).
 */

import React from 'react';
import {createNativeStackNavigator} from '@react-navigation/native-stack';
import {
  HomeScreen,
  ChaptersScreen,
  TechniquesScreen,
  TechniqueDetailScreen,
} from '@/screens';
import type {MainStackParamList} from './types';

const Stack = createNativeStackNavigator<MainStackParamList>();

const MainNavigator: React.FC = () => {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        contentStyle: {backgroundColor: '#f5f5f5'},
        animation: 'slide_from_right',
      }}>
      <Stack.Screen name="Home" component={HomeScreen} />
      <Stack.Screen name="Chapters" component={ChaptersScreen} />
      <Stack.Screen name="Techniques" component={TechniquesScreen} />
      <Stack.Screen name="TechniqueDetail" component={TechniqueDetailScreen} />
    </Stack.Navigator>
  );
};

export default MainNavigator;
