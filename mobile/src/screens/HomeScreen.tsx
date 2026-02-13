/**
 * Home Screen
 *
 * Displays a paginated grid of fighting books (e-commerce style).
 * Publicly accessible without authentication.
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
  Dimensions,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {useInfiniteQuery} from '@tanstack/react-query';
import {useNavigation} from '@react-navigation/native';
import type {NativeStackNavigationProp} from '@react-navigation/native-stack';
import {getFightingBooks} from '@/api/content';
import type {FightingBook} from '@/types/api';
import type {MainStackParamList} from '@/navigation/types';

type NavigationProp = NativeStackNavigationProp<MainStackParamList, 'Home'>;

const PAGE_SIZE = 20;
const COLUMN_COUNT = 2;
const CARD_GAP = 12;
const SCREEN_PADDING = 16;
const cardWidth =
  (Dimensions.get('window').width - SCREEN_PADDING * 2 - CARD_GAP) /
  COLUMN_COUNT;

const HomeScreen: React.FC = () => {
  const navigation = useNavigation<NavigationProp>();

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    isError,
    error,
    refetch,
    isRefetching,
  } = useInfiniteQuery({
    queryKey: ['fightingBooks'],
    queryFn: ({pageParam = 1}) =>
      getFightingBooks({page: pageParam, page_size: PAGE_SIZE}),
    getNextPageParam: (lastPage) => {
      if (lastPage.page < lastPage.total_pages) {
        return lastPage.page + 1;
      }
      return undefined;
    },
    initialPageParam: 1,
  });

  const books = data?.pages.flatMap((page) => page.data) ?? [];

  const handleBookPress = useCallback(
    (book: FightingBook) => {
      navigation.navigate('Chapters', {
        bookId: book.id,
        bookTitle: book.title,
      });
    },
    [navigation],
  );

  const handleEndReached = useCallback(() => {
    if (hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  const renderBookCard = useCallback(
    ({item}: {item: FightingBook}) => (
      <TouchableOpacity
        style={styles.bookCard}
        onPress={() => handleBookPress(item)}
        activeOpacity={0.7}
        testID={`book-card-${item.id}`}>
        <View style={styles.coverPlaceholder}>
          <Text style={styles.coverIcon}>&#9876;</Text>
        </View>
        <View style={styles.bookInfo}>
          <Text style={styles.bookTitle} numberOfLines={2}>
            {item.title}
          </Text>
          <Text style={styles.bookAuthor} numberOfLines={1}>
            {item.sword_master_name}
          </Text>
          {item.publication_year ? (
            <Text style={styles.bookYear}>{item.publication_year}</Text>
          ) : null}
        </View>
      </TouchableOpacity>
    ),
    [handleBookPress],
  );

  const renderFooter = useCallback(() => {
    if (!isFetchingNextPage) {
      return null;
    }
    return (
      <View style={styles.footerLoader}>
        <ActivityIndicator size="small" color="#1a1a2e" />
      </View>
    );
  }, [isFetchingNextPage]);

  const renderEmpty = useCallback(() => {
    if (isLoading) {
      return null;
    }
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyText}>No fighting books available yet.</Text>
      </View>
    );
  }, [isLoading]);

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <View>
          <Text style={styles.title}>HEMA Lessons</Text>
          <Text style={styles.subtitle}>Historical European Martial Arts</Text>
        </View>
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
              : 'Failed to load fighting books'}
          </Text>
          <TouchableOpacity style={styles.retryButton} onPress={() => refetch()}>
            <Text style={styles.retryText}>Retry</Text>
          </TouchableOpacity>
        </View>
      ) : (
        <FlatList
          data={books}
          renderItem={renderBookCard}
          keyExtractor={(item) => item.id.toString()}
          numColumns={COLUMN_COUNT}
          contentContainerStyle={styles.listContent}
          columnWrapperStyle={styles.columnWrapper}
          onEndReached={handleEndReached}
          onEndReachedThreshold={0.5}
          ListFooterComponent={renderFooter}
          ListEmptyComponent={renderEmpty}
          refreshControl={
            <RefreshControl
              refreshing={isRefetching && !isFetchingNextPage}
              onRefresh={refetch}
              colors={['#1a1a2e']}
              tintColor="#1a1a2e"
            />
          }
          testID="fighting-books-list"
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
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: SCREEN_PADDING,
    paddingTop: 12,
    paddingBottom: 16,
  },
  title: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#1a1a2e',
  },
  subtitle: {
    fontSize: 14,
    color: '#666',
    marginTop: 2,
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
    paddingHorizontal: SCREEN_PADDING,
    paddingBottom: 24,
  },
  columnWrapper: {
    justifyContent: 'space-between',
    marginBottom: CARD_GAP,
  },
  bookCard: {
    width: cardWidth,
    backgroundColor: '#fff',
    borderRadius: 16,
    overflow: 'hidden',
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 2},
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 4,
  },
  coverPlaceholder: {
    width: '100%',
    height: cardWidth * 0.75,
    backgroundColor: '#1a1a2e',
    justifyContent: 'center',
    alignItems: 'center',
  },
  coverIcon: {
    fontSize: 48,
    color: 'rgba(255,255,255,0.3)',
  },
  bookInfo: {
    padding: 12,
  },
  bookTitle: {
    fontSize: 15,
    fontWeight: '700',
    color: '#1a1a2e',
    marginBottom: 4,
  },
  bookAuthor: {
    fontSize: 13,
    color: '#666',
    marginBottom: 2,
  },
  bookYear: {
    fontSize: 12,
    color: '#999',
  },
  footerLoader: {
    paddingVertical: 20,
    alignItems: 'center',
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

export default HomeScreen;
