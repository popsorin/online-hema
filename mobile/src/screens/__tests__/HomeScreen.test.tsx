import React from 'react';
import {render, fireEvent, waitFor} from '@testing-library/react-native';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import HomeScreen from '../HomeScreen';

const mockNavigate = jest.fn();
jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({navigate: mockNavigate}),
}));

const mockGetFightingBooks = jest.fn();
jest.mock('@/api/content', () => ({
  getFightingBooks: (...args: unknown[]) => mockGetFightingBooks(...args),
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

describe('HomeScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the header title', async () => {
    mockGetFightingBooks.mockResolvedValue({
      data: [],
      page: 1,
      page_size: 20,
      total_count: 0,
      total_pages: 0,
    });

    const {getByText} = render(<HomeScreen />, {wrapper: createWrapper()});

    expect(getByText('HEMA Lessons')).toBeTruthy();
  });

  it('renders fighting book cards', async () => {
    mockGetFightingBooks.mockResolvedValue({
      data: [
        {
          id: 1,
          sword_master_id: 2,
          title: 'Fior di Battaglia',
          description: 'A combat manual',
          publication_year: 1409,
          cover_image_url: null,
          created_at: '2026-01-18T10:00:00Z',
          updated_at: '2026-01-18T10:00:00Z',
          sword_master_name: 'Fiore dei Liberi',
        },
      ],
      page: 1,
      page_size: 20,
      total_count: 1,
      total_pages: 1,
    });

    const {getByText} = render(<HomeScreen />, {wrapper: createWrapper()});

    await waitFor(() => {
      expect(getByText('Fior di Battaglia')).toBeTruthy();
    });
    expect(getByText('Fiore dei Liberi')).toBeTruthy();
    expect(getByText('1409')).toBeTruthy();
  });

  it('navigates to Chapters when a book card is pressed', async () => {
    mockGetFightingBooks.mockResolvedValue({
      data: [
        {
          id: 1,
          sword_master_id: 2,
          title: 'Fior di Battaglia',
          description: 'A combat manual',
          publication_year: 1409,
          cover_image_url: null,
          created_at: '2026-01-18T10:00:00Z',
          updated_at: '2026-01-18T10:00:00Z',
          sword_master_name: 'Fiore dei Liberi',
        },
      ],
      page: 1,
      page_size: 20,
      total_count: 1,
      total_pages: 1,
    });

    const {getByTestId} = render(<HomeScreen />, {wrapper: createWrapper()});

    await waitFor(() => {
      expect(getByTestId('book-card-1')).toBeTruthy();
    });

    fireEvent.press(getByTestId('book-card-1'));
    expect(mockNavigate).toHaveBeenCalledWith('Chapters', {
      bookId: 1,
      bookTitle: 'Fior di Battaglia',
    });
  });

  it('shows empty message when no books', async () => {
    mockGetFightingBooks.mockResolvedValue({
      data: [],
      page: 1,
      page_size: 20,
      total_count: 0,
      total_pages: 0,
    });

    const {getByText} = render(<HomeScreen />, {wrapper: createWrapper()});

    await waitFor(() => {
      expect(getByText('No fighting books available yet.')).toBeTruthy();
    });
  });
});
