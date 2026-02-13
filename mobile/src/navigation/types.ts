/**
 * Navigation Type Definitions
 *
 * Defines the navigation structure and screen params for type-safe navigation.
 */

import type {NativeStackScreenProps} from '@react-navigation/native-stack';
import type {Technique} from '@/types/api';

/**
 * Main Stack - all app screens (content is publicly accessible)
 */
export type MainStackParamList = {
  Home: undefined;
  Chapters: {bookId: number; bookTitle: string};
  Techniques: {chapterId: number; chapterTitle: string};
  TechniqueDetail: {technique: Technique};
};

export type MainStackScreenProps<T extends keyof MainStackParamList> =
  NativeStackScreenProps<MainStackParamList, T>;

// Declare global namespace for React Navigation
declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace ReactNavigation {
    interface RootParamList extends MainStackParamList {}
  }
}
