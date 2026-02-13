const {getDefaultConfig, mergeConfig} = require('@react-native/metro-config');
const path = require('path');

/**
 * Metro configuration
 * https://reactnative.dev/docs/metro
 *
 * @type {import('metro-config').MetroConfig}
 */
const defaultConfig = getDefaultConfig(__dirname);

const config = {
  resolver: {
    resolverMainFields: ['react-native', 'browser', 'main'],
    resolveRequest: (context, moduleName, platform) => {
      // Force axios to use the browser/RN compatible version
      if (moduleName === 'axios') {
        return {
          filePath: path.resolve(__dirname, 'node_modules/axios/dist/browser/axios.cjs'),
          type: 'sourceFile',
        };
      }
      // Use default resolution for everything else
      return context.resolveRequest(context, moduleName, platform);
    },
  },
};

module.exports = mergeConfig(defaultConfig, config);
