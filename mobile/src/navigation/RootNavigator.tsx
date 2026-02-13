/**
 * Root Navigator
 *
 * Top-level navigator that wraps the main content stack.
 */

import React from 'react';
import {NavigationContainer} from '@react-navigation/native';
import MainNavigator from './MainNavigator';

const RootNavigator: React.FC = () => {
  return (
    <NavigationContainer>
      <MainNavigator />
    </NavigationContainer>
  );
};

export default RootNavigator;
