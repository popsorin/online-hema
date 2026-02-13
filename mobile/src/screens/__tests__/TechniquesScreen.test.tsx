import React from 'react';
import {render, fireEvent, waitFor} from '@testing-library/react-native';
import {QueryClient, QueryClientProvider} from '@tanstack/react-query';
import TechniquesScreen from '../TechniquesScreen';

const mockNavigate = jest.fn();
const mockGoBack = jest.fn();
jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({navigate: mockNavigate, goBack: mockGoBack}),
  useRoute: () => ({
    params: {chapterId: 3, chapterTitle: 'Longsword'},
  }),
}));

const mockGetTechniques = jest.fn();
jest.mock('@/api/content', () => ({
  getTechniques: (...args: unknown[]) => mockGetTechniques(...args),
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

describe('TechniquesScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders the chapter title in the header', async () => {
    mockGetTechniques.mockResolvedValue([]);

    const {getByText} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    expect(getByText('Longsword')).toBeTruthy();
    expect(getByText('Techniques')).toBeTruthy();
  });

  it('renders technique list', async () => {
    mockGetTechniques.mockResolvedValue([
      {
        id: 1,
        chapter_id: 3,
        name: 'Posta di Donna',
        description: "The Woman's Guard",
        instructions: 'Hold the sword near your right shoulder',
        video_url: null,
        thumbnail_url: null,
        order_in_chapter: 1,
        created_at: '2026-01-18T10:00:00Z',
        updated_at: '2026-01-18T10:00:00Z',
      },
      {
        id: 2,
        chapter_id: 3,
        name: 'Zornhau',
        description: 'The Wrath Strike',
        instructions: 'Strike diagonally',
        video_url: null,
        thumbnail_url: null,
        order_in_chapter: 2,
        created_at: '2026-01-18T10:00:00Z',
        updated_at: '2026-01-18T10:00:00Z',
      },
    ]);

    const {getByText} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByText('Posta di Donna')).toBeTruthy();
    });
    expect(getByText('Zornhau')).toBeTruthy();
  });

  it('navigates to TechniqueDetail when a technique is pressed', async () => {
    const technique = {
      id: 1,
      chapter_id: 3,
      name: 'Posta di Donna',
      description: "The Woman's Guard",
      instructions: 'Hold the sword near your right shoulder',
      video_url: null,
      thumbnail_url: null,
      order_in_chapter: 1,
      created_at: '2026-01-18T10:00:00Z',
      updated_at: '2026-01-18T10:00:00Z',
    };
    mockGetTechniques.mockResolvedValue([technique]);

    const {getByTestId} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByTestId('technique-button-1')).toBeTruthy();
    });

    fireEvent.press(getByTestId('technique-button-1'));
    expect(mockNavigate).toHaveBeenCalledWith('TechniqueDetail', {
      technique,
    });
  });

  it('navigates back when back button is pressed', async () => {
    mockGetTechniques.mockResolvedValue([]);

    const {getByTestId} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    fireEvent.press(getByTestId('back-button'));
    expect(mockGoBack).toHaveBeenCalled();
  });

  it('shows empty message when no techniques', async () => {
    mockGetTechniques.mockResolvedValue([]);

    const {getByText} = render(<TechniquesScreen />, {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(getByText('No techniques available yet.')).toBeTruthy();
    });
  });
});
