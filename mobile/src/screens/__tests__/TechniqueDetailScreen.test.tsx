import React from 'react';
import {render, fireEvent} from '@testing-library/react-native';
import TechniqueDetailScreen from '../TechniqueDetailScreen';

const mockGoBack = jest.fn();
const mockRouteParams = {
  item: {
    id: 1,
    section_id: 3,
    kind: 'technique',
    title: 'Posta di Donna',
    description: "The Woman's Guard - a high guard position",
    position: 1,
    attributes: {
      instructions:
        "Hold the sword with the hilt near your right shoulder, point aimed at the opponent's face",
      historical_image_url: '',
    } as Record<string, string> | undefined,
  },
};

jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({goBack: mockGoBack}),
  useRoute: () => ({params: mockRouteParams}),
}));

describe('TechniqueDetailScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockRouteParams.item.attributes = {
      instructions:
        "Hold the sword with the hilt near your right shoulder, point aimed at the opponent's face",
      historical_image_url: '',
    };
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

  it('shows image placeholder when no historical_image_url', () => {
    mockRouteParams.item.attributes = undefined;

    const {getByTestId} = render(<TechniqueDetailScreen />);

    expect(getByTestId('image-placeholder')).toBeTruthy();
  });

  it('shows historical image when historical_image_url exists', () => {
    mockRouteParams.item.attributes = {
      instructions: 'Some instructions',
      historical_image_url:
        '/assets/books/fior-di-battaglia/techniques/posta-di-donna/historical.jpg',
    };

    const {getByTestId} = render(<TechniqueDetailScreen />);

    expect(getByTestId('historical-image')).toBeTruthy();
  });

  it('navigates back when back button is pressed', () => {
    const {getByTestId} = render(<TechniqueDetailScreen />);

    fireEvent.press(getByTestId('back-button'));
    expect(mockGoBack).toHaveBeenCalled();
  });
});
