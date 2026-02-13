/**
 * Chapters Screen
 *
 * Displays a list of chapters for a fighting book.
 * Each chapter is a clickable button that navigates to the techniques list.
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
import {getChapters} from '@/api/content';
import type {Chapter} from '@/types/api';
import type {MainStackParamList} from '@/navigation/types';

type NavigationProp = NativeStackNavigationProp<MainStackParamList, 'Chapters'>;
type ChaptersRouteProp = RouteProp<MainStackParamList, 'Chapters'>;

const ChaptersScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();
  const route = useRoute<ChaptersRouteProp>();
  const {bookId, bookTitle} = route.params;

  const {data: chapters, isLoading, isError, error, refetch, isRefetching} = useQuery({
    queryKey: ['chapters', bookId],
    queryFn: () => getChapters(bookId),
  });

  const handleChapterPress = useCallback(
    (chapter: Chapter) => {
      navigation.navigate('Techniques', {
        chapterId: chapter.id,
        chapterTitle: chapter.title,
      });
    },
    [navigation],
  );

  const renderChapter = useCallback(
    ({item}: {item: Chapter}) => (
      <TouchableOpacity
        style={styles.chapterButton}
        onPress={() => handleChapterPress(item)}
        activeOpacity={0.7}
        testID={`chapter-button-${item.id}`}>
        <View style={styles.chapterNumber}>
          <Text style={styles.chapterNumberText}>{item.chapter_number}</Text>
        </View>
        <View style={styles.chapterInfo}>
          <Text style={styles.chapterTitle}>{item.title}</Text>
          <Text style={styles.chapterDescription} numberOfLines={2}>
            {item.description}
          </Text>
        </View>
        <Text style={styles.chevron}>&#8250;</Text>
      </TouchableOpacity>
    ),
    [handleChapterPress],
  );

  const renderEmpty = useCallback(() => {
    if (isLoading) {
      return null;
    }
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyText}>No chapters available yet.</Text>
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
            {bookTitle}
          </Text>
          <Text style={styles.headerSubtitle}>Chapters</Text>
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
              : 'Failed to load chapters'}
          </Text>
          <TouchableOpacity style={styles.retryButton} onPress={() => refetch()}>
            <Text style={styles.retryText}>Retry</Text>
          </TouchableOpacity>
        </View>
      ) : (
        <FlatList
          data={chapters}
          renderItem={renderChapter}
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
          testID="chapters-list"
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
  chapterButton: {
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
  chapterNumber: {
    width: 44,
    height: 44,
    borderRadius: 12,
    backgroundColor: '#1a1a2e',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 14,
  },
  chapterNumberText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: 'bold',
  },
  chapterInfo: {
    flex: 1,
  },
  chapterTitle: {
    fontSize: 16,
    fontWeight: '700',
    color: '#1a1a2e',
    marginBottom: 4,
  },
  chapterDescription: {
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

export default ChaptersScreen;
