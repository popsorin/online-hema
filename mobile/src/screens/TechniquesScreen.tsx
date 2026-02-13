/**
 * Techniques Screen
 *
 * Displays a list of techniques for a chapter.
 * Each technique is a clickable button that navigates to the technique detail.
 */

import React, {useCallback} from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  FlatList,
  ActivityIndicator,
  RefreshControl,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {useQuery} from '@tanstack/react-query';
import {useNavigation, useRoute} from '@react-navigation/native';
import type {NativeStackNavigationProp} from '@react-navigation/native-stack';
import type {RouteProp} from '@react-navigation/native';
import {getTechniques} from '@/api/content';
import type {Technique} from '@/types/api';
import type {MainStackParamList} from '@/navigation/types';

type NavigationProp = NativeStackNavigationProp<
  MainStackParamList,
  'Techniques'
>;
type TechniquesRouteProp = RouteProp<MainStackParamList, 'Techniques'>;

const TechniquesScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();
  const route = useRoute<TechniquesRouteProp>();
  const {chapterId, chapterTitle} = route.params;

  const {
    data: techniques,
    isLoading,
    isError,
    error,
    refetch,
    isRefetching,
  } = useQuery({
    queryKey: ['techniques', chapterId],
    queryFn: () => getTechniques(chapterId),
  });

  const handleTechniquePress = useCallback(
    (technique: Technique) => {
      navigation.navigate('TechniqueDetail', {technique});
    },
    [navigation],
  );

  const renderTechnique = useCallback(
    ({item}: {item: Technique}) => (
      <TouchableOpacity
        style={styles.techniqueButton}
        onPress={() => handleTechniquePress(item)}
        activeOpacity={0.7}
        testID={`technique-button-${item.id}`}>
        <View style={styles.orderBadge}>
          <Text style={styles.orderText}>{item.order_in_chapter}</Text>
        </View>
        <View style={styles.techniqueInfo}>
          <Text style={styles.techniqueName}>{item.name}</Text>
          <Text style={styles.techniqueDescription} numberOfLines={2}>
            {item.description}
          </Text>
        </View>
        <Text style={styles.chevron}>&#8250;</Text>
      </TouchableOpacity>
    ),
    [handleTechniquePress],
  );

  const renderEmpty = useCallback(() => {
    if (isLoading) {
      return null;
    }
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyText}>No techniques available yet.</Text>
      </View>
    );
  }, [isLoading]);

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
            {chapterTitle}
          </Text>
          <Text style={styles.headerSubtitle}>Techniques</Text>
        </View>
        <View style={styles.headerSpacer} />
      </View>

      {isLoading ? (
        <View style={styles.loadingContainer}>
          <ActivityIndicator size="large" color="#1a1a2e" />
        </View>
      ) : isError ? (
        <View style={styles.errorContainer}>
          <Text style={styles.errorText}>
            {error instanceof Error
              ? error.message
              : 'Failed to load techniques'}
          </Text>
          <TouchableOpacity style={styles.retryButton} onPress={() => refetch()}>
            <Text style={styles.retryText}>Retry</Text>
          </TouchableOpacity>
        </View>
      ) : (
        <FlatList
          data={techniques}
          renderItem={renderTechnique}
          keyExtractor={(item) => item.id.toString()}
          contentContainerStyle={styles.listContent}
          ListEmptyComponent={renderEmpty}
          refreshControl={
            <RefreshControl
              refreshing={isRefetching}
              onRefresh={refetch}
              colors={['#1a1a2e']}
              tintColor="#1a1a2e"
            />
          }
          testID="techniques-list"
        />
      )}
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
  headerSubtitle: {
    fontSize: 13,
    color: '#666',
    marginTop: 2,
  },
  headerSpacer: {
    width: 40,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  errorText: {
    fontSize: 16,
    color: '#e53935',
    textAlign: 'center',
    marginBottom: 16,
  },
  retryButton: {
    backgroundColor: '#1a1a2e',
    borderRadius: 12,
    paddingHorizontal: 24,
    paddingVertical: 12,
  },
  retryText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: '600',
  },
  listContent: {
    paddingHorizontal: 16,
    paddingBottom: 24,
  },
  techniqueButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 16,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 2},
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 3,
  },
  orderBadge: {
    width: 36,
    height: 36,
    borderRadius: 18,
    backgroundColor: '#e8e8f0',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 14,
  },
  orderText: {
    color: '#1a1a2e',
    fontSize: 15,
    fontWeight: '700',
  },
  techniqueInfo: {
    flex: 1,
  },
  techniqueName: {
    fontSize: 16,
    fontWeight: '700',
    color: '#1a1a2e',
    marginBottom: 4,
  },
  techniqueDescription: {
    fontSize: 13,
    color: '#666',
    lineHeight: 18,
  },
  chevron: {
    fontSize: 24,
    color: '#ccc',
    marginLeft: 8,
  },
  emptyContainer: {
    paddingVertical: 60,
    alignItems: 'center',
  },
  emptyText: {
    fontSize: 16,
    color: '#999',
  },
});

export default TechniquesScreen;
