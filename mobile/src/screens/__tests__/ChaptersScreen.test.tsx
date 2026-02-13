import React from 'react';
import {render, fireEvent, waitFor} from '@testing-library/react-native';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import ChaptersScreen from '../ChaptersScreen';

const mockNavigate = jest.fn();
const mockGoBack = jest.fn();
jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({navigate: mockNavigate, goBack: mockGoBack}),
  useRoute: () => ({
    params: {bookId: 1, bookTitle: 'Fior di Battaglia'},
  }),
}));

const mockGetChapters = jest.fn();
jest.mock('@/api/content', () => ({
  getChapters: (...args: unknown[]) => mockGetChapters(...args),
}));

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {retry: false},
    },
  });
  return ({children}: {children: React.ReactNode}) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('ChaptersScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the book title in the header', async () => {
    mockGetChapters.mockResolvedValue([]);

    const {getByText} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    expect(getByText('Fior di Battaglia')).toBeTruthy();
    expect(getByText('Chapters')).toBeTruthy();
  });

  it('renders chapter list', async () => {
    mockGetChapters.mockResolvedValue([
      {
        id: 1,
        fighting_book_id: 1,
        chapter_number: 1,
        title: 'Wrestling',
        description: 'Techniques for unarmed combat and grappling',
        created_at: '2026-01-18T10:00:00Z',
        updated_at: '2026-01-18T10:00:00Z',
      },
      {
        id: 2,
        fighting_book_id: 1,
        chapter_number: 2,
        title: 'Dagger Combat',
        description: 'Fighting with the dagger in various situations',
        created_at: '2026-01-18T10:00:00Z',
        updated_at: '2026-01-18T10:00:00Z',
      },
    ]);

    const {getByText} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByText('Wrestling')).toBeTruthy();
    });
    expect(getByText('Dagger Combat')).toBeTruthy();
  });

  it('navigates to Techniques when a chapter is pressed', async () => {
    mockGetChapters.mockResolvedValue([
      {
        id: 3,
        fighting_book_id: 1,
        chapter_number: 3,
        title: 'Longsword',
        description: 'The art of fighting with the longsword',
        created_at: '2026-01-18T10:00:00Z',
        updated_at: '2026-01-18T10:00:00Z',
      },
    ]);

    const {getByTestId} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByTestId('chapter-button-3')).toBeTruthy();
    });

    fireEvent.press(getByTestId('chapter-button-3'));
    expect(mockNavigate).toHaveBeenCalledWith('Techniques', {
      chapterId: 3,
      chapterTitle: 'Longsword',
    });
  });

  it('navigates back when back button is pressed', async () => {
    mockGetChapters.mockResolvedValue([]);

    const {getByTestId} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    fireEvent.press(getByTestId('back-button'));
    expect(mockGoBack).toHaveBeenCalled();
  });

  it('shows empty message when no chapters', async () => {
    mockGetChapters.mockResolvedValue([]);

    const {getByText} = render(<ChaptersScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByText('No chapters available yet.')).toBeTruthy();
    });
  });
});
