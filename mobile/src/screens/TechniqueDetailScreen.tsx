/**
 * Technique Detail Screen
 *
 * Displays full details of an item: title, description, instructions,
 * and a historical image (or placeholder if none is available).
 */

import React from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {useNavigation, useRoute} from '@react-navigation/native';
import type {NativeStackNavigationProp} from '@react-navigation/native-stack';
import type {RouteProp} from '@react-navigation/native';
import type {MainStackParamList} from '@/navigation/types';

type NavigationProp = NativeStackNavigationProp<
  MainStackParamList,
  'TechniqueDetail'
>;
type DetailRouteProp = RouteProp<MainStackParamList, 'TechniqueDetail'>;

const TechniqueDetailScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();
  const route = useRoute<DetailRouteProp>();
  const {item} = route.params;

  const instructions = item.attributes?.instructions;
  const historicalImageUrl = item.attributes?.historical_image_url;

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity
          style={styles.backButton}
          onPress={() => navigation.goBack()}
          testID="back-button">
          <Text style={styles.backText}>&#8249;</Text>
        </TouchableOpacity>
        <View style={styles.headerTitleContainer}>
          <Text style={styles.headerTitle} numberOfLines={1}>
            Technique
          </Text>
        </View>
        <View style={styles.headerSpacer} />
      </View>

      <ScrollView
        contentContainerStyle={styles.scrollContent}
        testID="technique-detail-scroll">
        <Text style={styles.techniqueName}>{item.title}</Text>

        <View style={styles.imageSection}>
          {historicalImageUrl ? (
            <View style={styles.imagePlaceholder} testID="historical-image">
              <Text style={styles.placeholderIcon}>&#128247;</Text>
              <Text style={styles.placeholderText}>{historicalImageUrl}</Text>
            </View>
          ) : (
            <View style={styles.imagePlaceholder} testID="image-placeholder">
              <Text style={styles.placeholderIcon}>&#127909;</Text>
              <Text style={styles.placeholderText}>
                No image available yet
              </Text>
            </View>
          )}
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Description</Text>
          <Text style={styles.sectionText}>{item.description}</Text>
        </View>

        {instructions ? (
          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Instructions</Text>
            <Text style={styles.sectionText}>{instructions}</Text>
          </View>
        ) : null}
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingTop: 8,
    paddingBottom: 16,
  },
  backButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: '#fff',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 1},
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  backText: {
    fontSize: 28,
    color: '#1a1a2e',
    marginTop: -2,
  },
  headerTitleContainer: {
    flex: 1,
    alignItems: 'center',
    paddingHorizontal: 12,
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#1a1a2e',
  },
  headerSpacer: {
    width: 40,
  },
  scrollContent: {
    paddingHorizontal: 16,
    paddingBottom: 40,
  },
  techniqueName: {
    fontSize: 26,
    fontWeight: 'bold',
    color: '#1a1a2e',
    marginBottom: 20,
  },
  imageSection: {
    marginBottom: 24,
  },
  imagePlaceholder: {
    width: '100%',
    height: 200,
    borderRadius: 16,
    backgroundColor: '#e8e8f0',
    justifyContent: 'center',
    alignItems: 'center',
    borderWidth: 2,
    borderColor: '#d0d0d8',
    borderStyle: 'dashed',
  },
  placeholderIcon: {
    fontSize: 48,
    marginBottom: 8,
  },
  placeholderText: {
    fontSize: 15,
    color: '#999',
    fontWeight: '500',
    textAlign: 'center',
    paddingHorizontal: 16,
  },
  section: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 20,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 2},
    shadowOpacity: 0.06,
    shadowRadius: 8,
    elevation: 2,
  },
  sectionTitle: {
    fontSize: 14,
    fontWeight: '700',
    color: '#999',
    textTransform: 'uppercase',
    letterSpacing: 1,
    marginBottom: 10,
  },
  sectionText: {
    fontSize: 16,
    color: '#333',
    lineHeight: 24,
  },
});

export default TechniqueDetailScreen;
