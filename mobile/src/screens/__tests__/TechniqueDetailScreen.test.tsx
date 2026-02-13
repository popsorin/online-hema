import React from 'react';
import {render, fireEvent} from '@testing-library/react-native';
import {Linking} from 'react-native';
import TechniqueDetailScreen from '../TechniqueDetailScreen';

const mockGoBack = jest.fn();
const mockRouteParams = {
  technique: {
    id: 1,
    chapter_id: 3,
    name: 'Posta di Donna',
    description: "The Woman's Guard - a high guard position",
    instructions:
      'Hold the sword with the hilt near your right shoulder, point aimed at the opponent\'s face',
    video_url: null as string | null,
    thumbnail_url: null,
    order_in_chapter: 1,
    created_at: '2026-01-18T10:00:00Z',
    updated_at: '2026-01-18T10:00:00Z',
  },
};

jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({goBack: mockGoBack}),
  useRoute: () => ({params: mockRouteParams}),
}));

describe('TechniqueDetailScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockRouteParams.technique.video_url = null;
  });

  it('renders technique name', () => {
    const {getByText} = render(<TechniqueDetailScreen />);

    expect(getByText('Posta di Donna')).toBeTruthy();
  });

  it('renders description section', () => {
    const {getByText} = render(<TechniqueDetailScreen />);

    expect(getByText('Description')).toBeTruthy();
    expect(
      getByText("The Woman's Guard - a high guard position"),
    ).toBeTruthy();
  });

  it('renders instructions section', () => {
    const {getByText} = render(<TechniqueDetailScreen />);

    expect(getByText('Instructions')).toBeTruthy();
    expect(
      getByText(
        "Hold the sword with the hilt near your right shoulder, point aimed at the opponent's face",
      ),
    ).toBeTruthy();
  });

  it('shows video placeholder when no video_url', () => {
    const {getByTestId} = render(<TechniqueDetailScreen />);

    expect(getByTestId('video-placeholder')).toBeTruthy();
  });

  it('shows play button when video_url exists', () => {
    mockRouteParams.technique.video_url = 'https://example.com/video.mp4';

    const {getByTestId} = render(<TechniqueDetailScreen />);

    expect(getByTestId('video-button')).toBeTruthy();
  });

  it('opens video URL when play button pressed', () => {
    mockRouteParams.technique.video_url = 'https://example.com/video.mp4';
    jest.spyOn(Linking, 'openURL').mockResolvedValue(undefined as never);

    const {getByTestId} = render(<TechniqueDetailScreen />);

    fireEvent.press(getByTestId('video-button'));
    expect(Linking.openURL).toHaveBeenCalledWith(
      'https://example.com/video.mp4',
    );
  });

  it('navigates back when back button is pressed', () => {
    const {getByTestId} = render(<TechniqueDetailScreen />);

    fireEvent.press(getByTestId('back-button'));
    expect(mockGoBack).toHaveBeenCalled();
  });
});
